package signer

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
)

type TransactionClient struct {
	client *ClientWithResponses
}

func NewTransactionClient(server string, opts ...ClientOption) (*TransactionClient, error) {
	tc := new(TransactionClient)

	logrus.Infof("New transaction client")

	client, err := NewClientWithResponses(server, opts...)
	if err != nil {
		return nil, err
	}

	tc.client = client

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

		return tx.WithSignature(signer, *response.JSON200.SignedData)
	}
}
