package tlsconfig

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
)

type TLSCertConfig struct {
	ClientKeyFile      string `json:"clientKeyFile"`
	ClientCertFile     string `json:"clientCertFile"`
	CACertFile         string `json:"caCertFile"`
	InsecureSkipVerify bool   `json:"insecureSkipVerify"`
}

func GetTLSConfig(config *TLSCertConfig) (*tls.Config, error) {
	if config == nil {
		return &tls.Config{}, nil
	}

	result := &tls.Config{}

	cert, err := tls.LoadX509KeyPair(config.ClientCertFile, config.ClientKeyFile)
	if err != nil {
		return result, fmt.Errorf("Error creating x509 keypair from client cert file %s and client key file %s", config.ClientCertFile, config.ClientKeyFile)
	}

	result.Certificates = []tls.Certificate{cert}

	caCert, err := ioutil.ReadFile(config.CACertFile)
	if err != nil {
		return result, fmt.Errorf("Error opening cert file %s, Error: %s", config.CACertFile, err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	result.RootCAs = caCertPool

	result.InsecureSkipVerify = config.InsecureSkipVerify

	return result, nil
}
