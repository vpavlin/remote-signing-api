package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
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
	handlers.SetupSigNonce(e, config)

	e.GET("/api/v1/", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, e.Routes())
	})

	if len(config.Server.CertPath) > 0 && len(config.Server.KeyPath) > 0 {
		logrus.Infof("Starting TLS server at %s", fmt.Sprintf("https://%s:%d", config.Server.Hostname, config.Server.Port))
		s := http.Server{
			Addr:      fmt.Sprintf("%s:%d", config.Server.Hostname, config.Server.Port),
			Handler:   e,
			TLSConfig: getTLSConfig(config.Server),
		}
		logrus.Fatal(s.ListenAndServeTLS(config.Server.CertPath, config.Server.KeyPath))
	}

	logrus.Fatal(e.Start(fmt.Sprintf("%s:%d", config.Server.Hostname, config.Server.Port)))
}

func getTLSConfig(config *config.Server) *tls.Config {
	var caCert []byte
	var err error
	var caCertPool *x509.CertPool

	caCert, err = ioutil.ReadFile(config.CACertPath)
	if err != nil {
		log.Fatal("Error opening cert file", config.CACertPath, ", error ", err)
	}
	caCertPool = x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	return &tls.Config{
		ServerName: fmt.Sprintf("%s:%d", config.Hostname, config.Port),
		// ClientAuth: tls.NoClientCert,				// Client certificate will not be requested and it is not required
		// ClientAuth: tls.RequestClientCert,			// Client certificate will be requested, but it is not required
		// ClientAuth: tls.RequireAnyClientCert,		// Client certificate is required, but any client certificate is acceptable
		// ClientAuth: tls.VerifyClientCertIfGiven,		// Client certificate will be requested and if present must be in the server's Certificate Pool
		// ClientAuth: tls.RequireAndVerifyClientCert,	// Client certificate will be required and must be present in the server's Certificate Pool
		ClientAuth: 4,
		ClientCAs:  caCertPool,
		MinVersion: tls.VersionTLS12, // TLS versions below 1.2 are considered insecure - see https://www.rfc-editor.org/rfc/rfc7525.txt for details
	}
}
