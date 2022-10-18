package main

import (
	"log"

	"github.com/jessevdk/go-flags"
	"github.com/vpavlin/remote-signing-api/pkg/nonce"
)

type NonceArgs struct {
	ChainId       uint64 `long:"chain" description:"CHain Id" required:"true"`
	Address       string `short:"a" description:"Public key of the signer" required:"true"`
	Key           string `short:"k" description:"API Key for the signing service" required:"true"`
	Action        string `long:"action" description:"Action to perform" required:"true"`
	Server        string `long:"server" description:"Signer Server URL" default:"https://localhost:4444"`
	SkipTlsVerify bool   `long:"skipVerify" description:"Skip TLS Verify"`
}

var opts NonceArgs

func main() {

	var err error

	parser := flags.NewParser(&opts, flags.Default)

	_, err = parser.Parse()
	if err != nil {
		log.Fatal(err)
	}

	switch opts.Action {
	case "getWithSigner":
		err = getNonceWithSigner(opts.ChainId, opts.Address, opts.Key)
	case "get":
		err = getNonce(opts.ChainId, opts.Address, opts.Key)
	}

	if err != nil {
		log.Fatal(err)
	}

}

func getNonceWithSigner(chainId uint64, address string, apiKey string) error {
	client, err := nonce.NewNonceClientWithSigner(opts.Server, opts.SkipTlsVerify)
	if err != nil {
		return err
	}

	nonce, nonceReturnFN, err := client.GetNonceWithSigner(chainId, address, apiKey)
	if err != nil {
		nonceReturnFN()
		return err
	}

	log.Println("Got nonce: ", nonce)

	//nonceReturnFN()

	return nil
}

func getNonce(chainId uint64, address string, apiKey string) error {
	client, err := nonce.NewNonceClient(opts.Server, opts.SkipTlsVerify, opts.Key)
	if err != nil {
		return err
	}

	nonce, nonceReturnFN, err := client.GetNonce(chainId, address)
	if err != nil {
		if nonceReturnFN != nil {
			nonceReturnFN()
		}
		return err
	}

	log.Println("Got nonce: ", nonce)

	return nil
}
