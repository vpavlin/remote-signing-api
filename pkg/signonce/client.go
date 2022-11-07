package signonce

import (
	"context"
	"fmt"
	"net/http"

	tlsconfig "github.com/vpavlin/remote-signing-api/pkg/tlsConfig"
)

type NonceClient struct {
	nonceClient *ClientWithResponses
	apiKey      string
}

func NewNonceClient(server string, config *tlsconfig.TLSCertConfig, apiKey string) (*NonceClient, error) {
	clientOpt := func(c *Client) error {
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
	}

	nonceClient, err := NewClientWithResponses(server, clientOpt)
	if err != nil {
		return nil, err
	}

	return &NonceClient{
		nonceClient: nonceClient,
		apiKey:      apiKey,
	}, nil
}

func (c *NonceClient) GetNonce(contract string, chainId uint64, address string) (uint64, func() error, error) {
	responseNonce, err := c.nonceClient.GetNonceWithResponse(context.Background(), contract, chainId, address, c.addAuth)
	if err != nil {
		return 0, nil, err
	}

	if responseNonce.StatusCode() != 200 {
		return 0, nil, fmt.Errorf("Failed to get nonce for %s (%d) - code: %d", address, chainId, responseNonce.StatusCode())
	}

	fn := func() error {
		nonceToReturn := *responseNonce.JSON200.Nonce

		response, err := c.nonceClient.ReturnNonce(context.Background(), contract, chainId, address, nonceToReturn, c.addAuth)
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

func (c *NonceClient) ReturnNonce(contract string, chainId uint64, address string, nonce uint64) error {
	response, err := c.nonceClient.ReturnNonce(context.Background(), contract, chainId, address, nonce, c.addAuth)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to return nonce %d for %d - %s with code %d", nonce, chainId, address, response.StatusCode)
	}

	return nil
}

func (c *NonceClient) addAuth(ctx context.Context, req *http.Request) error {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	return nil
}
