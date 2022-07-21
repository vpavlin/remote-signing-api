package main

import (
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/vpavlin/remote-signing-api/handlers"
	"github.com/vpavlin/remote-signing-api/internal/nonce"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())

	client, err := ethclient.Dial("https://matic-testnet-archive-rpc.bwarelabs.com")
	if err != nil {
		log.Fatal(err)
	}
	nm, err := nonce.NewNonceManager(client)
	if err != nil {
		log.Fatal(err)
	}

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("NonceManager", nm)
			return next(c)
		}
	})

	//e.POST("/sign/:chainId/:address", handlers.HandleSign)
	e.GET("/nonce/:chainId/:address", handlers.HandleGetNonce)
	e.PUT("/nonce/:chainId/:address/:nonce", handlers.HandleReturnNonce)
	e.DELETE("/nonce/:chainId/:address", handlers.HandleDecreaseNonce)
	e.POST("/nonce/:chainId/:address/sync", handlers.HandleSync)

	e.Start("localhost:4444")

}
