package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"github.com/vpavlin/remote-signing-api/config"
	"github.com/vpavlin/remote-signing-api/internal/nonce"
	"github.com/vpavlin/remote-signing-api/pkg/signer"
	signonceServer "github.com/vpavlin/remote-signing-api/pkg/signonce"
)

type SigNonceHandler struct{}

func SetupSigNonce(e *echo.Echo, config *config.Config) {
	g := e.Group("")

	nm, err := nonce.NewNonceManager(config.RpcUrls, config.NonceManager)
	if err != nil {
		logrus.Fatal(err)
	}

	g.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("NonceManager", nm)
			return next(c)
		}
	})

	if config.NonceManager.AuthBySig {
		g.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
			KeyLookup: "header:X-NONCE-AUTH-HASH",
			Validator: func(hash string, ctx echo.Context) (bool, error) {
				address := ctx.Request().Header.Get("X-NONCE-AUTH-SIGNER")
				signature := ctx.Request().Header.Get("X-NONCE-AUTH-SIGNATURE")

				nonceAddress := ctx.Param("address")
				if len(nonceAddress) > 0 {
					if nonceAddress != address {
						return false, nil
					}
				}

				hashBytes, err := signer.StringToBytes(hash)
				if err != nil {
					return false, err
				}

				sigBytes, err := signer.StringToBytes(signature)
				if err != nil {
					return false, err
				}

				ok, err := signer.IsValidSignature(address, hashBytes, sigBytes)
				if err != nil {
					return false, err
				}

				if ok {
					return true, nil
				}

				return false, nil

			},
		}))
	} else {
		g.Use(middleware.KeyAuth(
			func(auth string, c echo.Context) (bool, error) {
				return config.NonceManager.ApiKey == auth, nil
			}))
	}

	nh := SigNonceHandler{}

	signonceServer.RegisterHandlers(g, &nh)
}

func (nh SigNonceHandler) GetNonceWithSigner(ctx echo.Context, contract string, chainId uint64, address string, params signonceServer.GetNonceWithSignerParams) error {
	nm := ctx.Get("NonceManager").(*nonce.NonceManager)

	resp := &signonceServer.NonceResponse{}

	nonce, err := nm.GetNonce(nonce.ChainID(chainId), address, &contract)
	if err != nil {
		ctx.Error(err)
		return ctx.JSON(http.StatusInternalServerError, &ErrorResponse{Error: err.Error()})
	}

	resp.Nonce = &nonce
	resp.Address = &address
	resp.ChainId = &chainId

	return ctx.JSON(http.StatusOK, resp)
}

func (nh SigNonceHandler) ReturnNonceWithSigner(ctx echo.Context, contract string, chainId uint64, address string, nonceI uint64, params signonceServer.ReturnNonceWithSignerParams) error {
	nm := ctx.Get("NonceManager").(*nonce.NonceManager)

	err := nm.ReturnNonce(nonceI, nonce.ChainID(chainId), address, &contract)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &ErrorResponse{Error: err.Error()})
	}

	return ctx.NoContent(http.StatusOK)
}

func (nh SigNonceHandler) GetNonce(ctx echo.Context, contract string, chainId uint64, address string) error {
	nm := ctx.Get("NonceManager").(*nonce.NonceManager)

	resp := &signonceServer.NonceResponse{}

	nonce, err := nm.GetNonce(nonce.ChainID(chainId), address, &contract)
	if err != nil {
		ctx.Error(err)
		return ctx.JSON(http.StatusInternalServerError, &ErrorResponse{Error: err.Error()})
	}

	resp.Nonce = &nonce
	resp.Address = &address
	resp.ChainId = &chainId

	return ctx.JSON(http.StatusOK, resp)
}

func (nh SigNonceHandler) ReturnNonce(ctx echo.Context, contract string, chainId uint64, address string, nonceI uint64) error {
	nm := ctx.Get("NonceManager").(*nonce.NonceManager)

	err := nm.ReturnNonce(nonceI, nonce.ChainID(chainId), address, &contract)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &ErrorResponse{Error: err.Error()})
	}

	return ctx.NoContent(http.StatusOK)
}

func (nh SigNonceHandler) SyncNonce(ctx echo.Context, contract string, chainId uint64, address string) error {
	nm := ctx.Get("NonceManager").(*nonce.NonceManager)

	_, err := nm.Sync(nonce.ChainID(chainId), address, &contract)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &ErrorResponse{Error: err.Error()})
	}

	return ctx.NoContent(http.StatusOK)
}
