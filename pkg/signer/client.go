package signer

import (
	"context"
	"fmt"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
	tlsconfig "github.com/vpavlin/remote-signing-api/pkg/tlsConfig"
)

type TransactionClient struct {
	client *ClientWithResponses
}

func NewSignerClientWithTLSOpts(server string, config *tlsconfig.TLSCertConfig, opts ...ClientOption) (*ClientWithResponses, error) {
	opts = append([]ClientOption{func(c *Client) error {
		if config != nil && c.Client == nil {
			tlsconf, err := tlsconfig.GetTLSConfig(config)
			if err != nil {
				return err
			}
			tr := &http.Transport{
				TLSClientConfig: tlsconf,
			}
			c.Client = &http.Client{Transport: tr}
		}

		return nil
	}}, opts...)

	client, err := NewClientWithResponses(server, opts...)

	return client, err
}

func NewTransactionClient(server string, config *tlsconfig.TLSCertConfig, opts ...ClientOption) (*TransactionClient, error) {
	tc := new(TransactionClient)

	logrus.Infof("New transaction client")

	opts = append([]ClientOption{func(c *Client) error {
		if config != nil && c.Client == nil {
			tlsconf, err := tlsconfig.GetTLSConfig(config)
			if err != nil {
				return err
			}
			tr := &http.Transport{
				TLSClientConfig: tlsconf,
			}
			c.Client = &http.Client{Transport: tr}
		}

		return nil
	}}, opts...)

	client, err := NewClientWithResponses(server, opts...)

	tc.client = client

	res, err := client.HealthWithResponse(context.Background())
	if err != nil {
		return nil, err
	}

	if res.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("Failed to initialize the signer client: Signer API not healthy: %s", res.Status())
	}

	return tc, nil
}

func (tc *TransactionClient) Signer(chainId *big.Int, apiKey string) func(address common.Address, tx *types.Transaction) (*types.Transaction, error) {
	logrus.Infof("New signer")
	return func(address common.Address, tx *types.Transaction) (*types.Transaction, error) {
		signer := types.LatestSignerForChainID(chainId)
		txBytes := signer.Hash(tx).Bytes()

		logrus.Infof("TX bytes: %s", txBytes)

		reqBody := new(SignBytes)
		reqBody.Bytes = &txBytes

		params := SignBytesParams{Authorization: "Bearer " + apiKey}

		response, err := tc.client.SignBytesWithResponse(context.Background(), address.String(), &params, *reqBody)
		if err != nil {
			return nil, err
		}

		if response.StatusCode() != 200 {
			return nil, fmt.Errorf("Failed to get signature: %s", response.Status())
		}

		return tx.WithSignature(signer, *response.JSON200.SignedData)
	}
}

func (tc *TransactionClient) SignBytes(address string, apiKey string, data []byte) ([]byte, error) {
	params := SignBytesParams{Authorization: "Bearer " + apiKey}

	reqBody := SignBytes{
		Bytes: &data,
	}

	response, err := tc.client.SignBytesWithResponse(context.Background(), address, &params, reqBody)
	if err != nil {
		return nil, err
	}

	if response.StatusCode() != 200 {
		return nil, fmt.Errorf("Failed to get signature: %s", response.Status())
	}

	return *response.JSON200.SignedData, nil
}
