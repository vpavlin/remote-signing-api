package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/crypto"
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
}

func main() {
	var opts SignCmd

	_, err := flags.Parse(&opts)
	if err != nil {
		panic(err)
	}

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
		log.Fatal("Status not 200: %s", rtx.Status())
	}

	ok, err := signer.IsValidSignature(opts.Address, hash, *rtx.JSON200.SignedData)
	if err != nil {
		log.Fatal(err)
	}

	if !ok {
		log.Fatal("Recovered address does not match")
	}
}

func produceHash(message string) []byte {
	prefix := []byte("\x19Ethereum Signed Message\n")
	length := []byte(strconv.Itoa(len(message)))
	prefixed := append(append(prefix[:], length[:]...), []byte(message)[:]...)
	hash := crypto.Keccak256Hash(prefixed)
	return hash.Bytes()
}
