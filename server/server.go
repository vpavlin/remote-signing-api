package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/vpavlin/remote-signing-api/config"
	"github.com/vpavlin/remote-signing-api/handlers"

	logrus "github.com/sirupsen/logrus"
)

func main() {
	e := echo.New()
	e.Use(
		middleware.Recover(),   // Recover from all panics to always have your server up
		middleware.Logger(),    // Log everything to stdout
		middleware.RequestID(), // Generate a request id on the HTTP response headers for identification
	)

	config, err := config.NewConfig(os.Args[1])
	if err != nil {
		logrus.Fatal(err)
	}

	ll, err := logrus.ParseLevel(config.Server.LogLevel)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.SetLevel(ll)

	handlers.SetupNonce(e, config)
	handlers.SetupSigner(e, config)

	e.GET("/api/v1/", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, e.Routes())
	})

	if len(config.Server.CertPath) > 0 && len(config.Server.KeyPath) > 0 {
		logrus.Fatal(e.StartTLS(fmt.Sprintf("%s:%d", config.Server.Hostname, config.Server.Port), config.Server.CertPath, config.Server.KeyPath))
	}

	logrus.Fatal(e.Start(fmt.Sprintf("%s:%d", config.Server.Hostname, config.Server.Port)))
}
