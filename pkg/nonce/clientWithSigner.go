package nonce

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/vpavlin/remote-signing-api/pkg/signer"
)

type ClientWithSigner struct {
	nonceClient  *ClientWithResponses
	signerClient *signer.ClientWithResponses
}

func NewNonceClientWithSigner(server string, insecureSkipVerify bool) (*ClientWithSigner, error) {
	client := new(ClientWithSigner)

	skipVerify := func(c *Client) error {
		if insecureSkipVerify && c.Client == nil {
			tr := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
			c.Client = &http.Client{Transport: tr}
		}

		return nil
	}

	nonceClient, err := NewClientWithResponses(server, skipVerify)
	if err != nil {
		return nil, err
	}

	client.nonceClient = nonceClient

	signerSkipVerify := func(c *signer.Client) error {
		if insecureSkipVerify && c.Client == nil {
			tr := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
			c.Client = &http.Client{Transport: tr}
		}

		return nil
	}

	signerClient, err := signer.NewClientWithResponses(server, signerSkipVerify)
	if err != nil {
		return nil, err
	}

	client.signerClient = signerClient

	return client, nil

}

func (c *ClientWithSigner) GetNonceWithSigner(chainId uint64, address string, apiKey string) (uint64, func() error, error) {
	hash := c.getHash(chainId, address)

	response, err := c.signerClient.SignBytesWithResponse(
		context.Background(),
		address,
		&signer.SignBytesParams{Authorization: fmt.Sprintf("Bearer %s", apiKey)},
		signer.SignBytesJSONRequestBody{Bytes: &hash},
	)

	if err != nil {
		return 0, nil, err
	}

	signature := response.JSON200.SignedData

	params := &GetNonceWithSignerParams{}
	params.XNONCEAUTHHASH = signer.BytesToString(hash)
	params.XNONCEAUTHSIGNATURE = signer.BytesToString(*signature)
	params.XNONCEAUTHSIGNER = address

	responseNonce, err := c.nonceClient.GetNonceWithSignerWithResponse(context.Background(), chainId, address, params)
	if err != nil {
		return 0, nil, err
	}

	fn := func() error {
		nonceToReturn := *responseNonce.JSON200.Nonce

		response, err := c.nonceClient.ReturnNonceWithSigner(context.Background(), chainId, address, nonceToReturn, (*ReturnNonceWithSignerParams)(params))
		if err != nil {
			return err
		}

		if response.StatusCode != http.StatusOK {
			return fmt.Errorf("Failed to return nonce %d for %d - %s with code %d", nonceToReturn, chainId, address, response.StatusCode)
		}

		return nil
	}

	return *responseNonce.JSON200.Nonce, fn, nil
}

func (c *ClientWithSigner) ReturnNonce(chainId uint64, address string, apiKey string, nonce uint64) error {
	hash := c.getHash(chainId, address)

	response, err := c.signerClient.SignBytesWithResponse(
		context.Background(),
		address,
		&signer.SignBytesParams{Authorization: fmt.Sprintf("Bearer %s", apiKey)},
		signer.SignBytesJSONRequestBody{Bytes: &hash},
	)

	if err != nil {
		return err
	}

	signature := response.JSON200.SignedData

	params := &GetNonceWithSignerParams{}
	params.XNONCEAUTHHASH = signer.BytesToString(hash)
	params.XNONCEAUTHSIGNATURE = signer.BytesToString(*signature)
	params.XNONCEAUTHSIGNER = address

	res, err := c.nonceClient.ReturnNonceWithSigner(context.Background(), chainId, address, nonce, (*ReturnNonceWithSignerParams)(params))
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to return nonce %d for %d - %s with code %d", nonce, chainId, address, res.StatusCode)
	}

	return nil

}

func (c *ClientWithSigner) getHash(chainId uint64, address string) []byte {
	message := fmt.Sprintf("%d-%s-%s-%d", chainId, address, "getNonce", time.Now().Unix())
	hash := crypto.Keccak256Hash([]byte(message))
	return hash.Bytes()
}
