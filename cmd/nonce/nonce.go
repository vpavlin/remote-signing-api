package main

import (
	"log"

	"github.com/jessevdk/go-flags"
	"github.com/vpavlin/remote-signing-api/pkg/nonce"
)

type NonceArgs struct {
	ChainId uint64 `long:"chain" description:"CHain Id" required:"true"`
	Address string `short:"a" description:"Public key of the signer" required:"true"`
	Key     string `short:"k" description:"API Key for the signing service" required:"true"`
	Action  string `long:"action" description:"Action to perform" required:"true"`
}

func main() {

	var opts NonceArgs
	var err error

	parser := flags.NewParser(&opts, flags.Default)

	_, err = parser.Parse()
	if err != nil {
		log.Fatal(err)
	}

	switch opts.Action {
	case "get":
		err = getNonce(opts.ChainId, opts.Address, opts.Key)
	}

	if err != nil {
		log.Fatal(err)
	}

}

func getNonce(chainId uint64, address string, apiKey string) error {
	client, err := nonce.NewNonceClientWithSigner("http://localhost:4444")
	if err != nil {
		return err
	}

	nonce, err := client.GetNonceWithSigner(chainId, address, apiKey)
	if err != nil {
		return err
	}

	log.Println("Got nonce: ", nonce)
	return nil
}
