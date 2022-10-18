package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"reflect"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/vpavlin/remote-signing-api/bindings"
	"github.com/vpavlin/remote-signing-api/pkg/nonce"
	"github.com/vpavlin/remote-signing-api/pkg/signer"

	"github.com/jessevdk/go-flags"
)

type ToSign struct {
	Address  string
	Message  string
	Deadline int
}

const URL = "http://localhost:4444/"

type SignCmd struct {
	Address string `short:"a" description:"Public key of the signer" required:"true"`
	Message string `short:"m" description:"Message to be signed" required:"true"`
	Key     string `short:"k" description:"API Key for the signing service" required:"true"`
	Action  string `long:"action" descroption:"Action to perform" required:"true" default:"byte"`
}

var opts SignCmd

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		panic(err)
	}

	switch opts.Action {
	case "byte":
		signBytes()
	case "tx":
		signTX()
	default:
		log.Fatal("Unknown action")
	}

}

func produceHash(message string) []byte {
	prefix := []byte("\x19Ethereum Signed Message\n")
	length := []byte(strconv.Itoa(len(message)))
	prefixed := append(append(prefix[:], length[:]...), []byte(message)[:]...)
	hash := crypto.Keccak256Hash(prefixed)
	return hash.Bytes()
}

func signBytes() {
	hash := produceHash(opts.Message)

	body := signer.SignBytesJSONBody{}
	body.Bytes = &hash

	params := signer.SignBytesParams{
		Authorization: fmt.Sprintf("Bearer %s", opts.Key),
	}

	c, err := signer.NewClientWithResponses(URL)
	if err != nil {
		log.Fatal(err)
	}

	rtx, err := c.SignBytesWithResponse(context.Background(), opts.Address, &params, body)
	if rtx.StatusCode() != http.StatusOK {
		log.Fatalf("Status not 200: %s", rtx.Status())
	}

	ok, err := signer.IsValidSignature(opts.Address, hash, *rtx.JSON200.SignedData)
	if err != nil {
		log.Fatal(err)
	}

	if !ok {
		log.Fatal("Recovered address does not match")
	}
}

func signTX() {

	userAddress := common.HexToAddress(opts.Address)
	contractAddress := common.HexToAddress("0x9b0e9c890b18babef972c48d4ae7939a52972e83")

	client, err := ethclient.Dial("https://matic-testnet-archive-rpc.bwarelabs.com")
	if err != nil {
		log.Fatal(err)
	}

	chainId, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	cm, err := bindings.NewCollectionManager(contractAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	sc, err := signer.NewTransactionClient(URL, true)
	if err != nil {
		log.Fatal(err)
	}

	callOpts := bind.CallOpts{}
	creator, err := cm.Creator(&callOpts)
	if err != nil {
		log.Fatal(err)
	}

	if userAddress != creator {
		//log.Fatalf("Only callable by the creator!")
	}

	nonceClient, err := nonce.NewNonceClientWithSigner(URL, true)
	if err != nil {
		log.Fatal(err)
	}

	nonce, returnNonceFN, err := nonceClient.GetNonceWithSigner(chainId.Uint64(), opts.Address, opts.Key)
	if err != nil {
		log.Fatal(err)
	}

	txOpts := &bind.TransactOpts{
		From:    userAddress,
		Signer:  sc.Signer(chainId, opts.Key),
		Context: context.Background(),
		Nonce:   big.NewInt(int64(nonce)),
		NoSend:  true,
	}

	tx, cerr := cm.CreateCollectionClone(txOpts)
	if cerr != nil {
		returnNonceFN()

		t := reflect.TypeOf(cerr)
		var data []byte

		for i := 0; i < t.Elem().NumField(); i++ {
			switch t.Elem().Field(i).Name {
			case "Data":
				data, err = hex.DecodeString(reflect.ValueOf(cerr).Elem().Field(i).Interface().(string)[2:])
				if err != nil {
					log.Fatal(err)
				}
			}
		}

		arguments := make(abi.Arguments, 2)
		arguments[0].Name = "caller"
		newType, err := abi.NewType("address", "address", []abi.ArgumentMarshaling{})
		if err != nil {
			log.Fatal(err)
		}
		arguments[0].Type = newType
		arguments[1].Name = "expected"
		arguments[1].Type = newType
		tryErr := abi.NewError("UnexpectedCaller", arguments)
		unpacked, err := tryErr.Unpack(data)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Unpacked UnexpectedCaller arguments: ", unpacked)

		log.Fatal("Call contract: ", cerr)
	}

	signer := types.LatestSignerForChainID(tx.ChainId())
	sender, err := types.Sender(signer, tx)
	if err != nil {

		log.Fatal(fmt.Errorf("Failed to check sender: %s", err))
	}

	err = client.SendTransaction(context.Background(), tx)
	if err != nil {
		returnNonceFN()
		log.Fatal(err)
	}

	log.Println(tx.Hash().Hex())

	log.Println(sender)

}

type jsonError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (err jsonError) Error() string {
	if err.Message == "" {
		return fmt.Sprintf("json-rpc error %d", err.Code)
	}
	return err.Message
}

func (err jsonError) ErrorCode() int {
	return err.Code
}

func (err jsonError) ErrorData() interface{} {
	return err.Data
}

type UnexpectedCaller struct {
	Caller   common.Address `json:"caller"`
	Expected common.Address `json:"expected"`
}
