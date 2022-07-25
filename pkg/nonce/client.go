package nonce

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/vpavlin/remote-signing-api/pkg/signer"
)

type ClientWithSigner struct {
	nonceClient  *ClientWithResponses
	signerClient *signer.ClientWithResponses
}

func NewNonceClientWithSigner(server string) (*ClientWithSigner, error) {
	client := new(ClientWithSigner)
	nonceClient, err := NewClientWithResponses(server)
	if err != nil {
		return nil, err
	}

	client.nonceClient = nonceClient

	signerClient, err := signer.NewClientWithResponses(server)
	if err != nil {
		return nil, err
	}

	client.signerClient = signerClient

	return client, nil

}

func (c *ClientWithSigner) GetNonceWithSigner(chainId uint64, address string, apiKey string) (uint64, error) {
	hash := c.getHash(chainId, address)

	response, err := c.signerClient.PostSignerAddressBytesWithResponse(
		context.Background(),
		address,
		&signer.PostSignerAddressBytesParams{Authorization: fmt.Sprintf("Bearer %s", apiKey)},
		signer.PostSignerAddressBytesJSONRequestBody{Bytes: &hash},
	)

	if err != nil {
		return 0, err
	}

	signature := response.JSON200.SignedData

	params := new(GetNonceParams)
	params.XNONCEAUTHHASH = signer.BytesToString(hash)
	params.XNONCEAUTHSIGNATURE = signer.BytesToString(*signature)
	params.XNONCEAUTHSIGNER = address

	responseNonce, err := c.nonceClient.GetNonceWithResponse(context.Background(), chainId, address, params)
	if err != nil {
		return 0, err
	}

	return *responseNonce.JSON200.Nonce, nil
}

func (c *ClientWithSigner) getHash(chainId uint64, address string) []byte {
	message := fmt.Sprintf("%d-%s-%s-%d", chainId, address, "getNonce", time.Now().Unix())
	hash := crypto.Keccak256Hash([]byte(message))
	return hash.Bytes()
}
