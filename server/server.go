package main

import (
	"fmt"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/vpavlin/remote-signing-api/config"
	"github.com/vpavlin/remote-signing-api/handlers"
	"github.com/vpavlin/remote-signing-api/internal/nonce"

	logrus "github.com/sirupsen/logrus"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())

	config, err := config.NewConfig(os.Args[1])
	if err != nil {
		logrus.Fatal(err)

	}

	ll, err := logrus.ParseLevel(config.Server.LogLevel)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.SetLevel(ll)

	nm, err := nonce.NewNonceManager(config.RpcUrls, config.NonceManager)
	if err != nil {
		logrus.Fatal(err)
	}

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("NonceManager", nm)
			return next(c)
		}
	})

	handlers.SeuptGroup(e)

	//e.POST("/sign/:chainId/:address", handlers.HandleSign)

	e.Start(fmt.Sprintf("%s:%d", config.Server.Hostname, config.Server.Port))
}
