package config

import (
	"encoding/json"
	"io/ioutil"
)

type Server struct {
	Hostname   string `json:"hostname"`
	Port       int    `json:"port"`
	LogLevel   string `json:"logLevel"`
	CertPath   string `json:"certPath"`
	KeyPath    string `json:"keyPath"`
	CACertPath string `json:"caCertPath"`
}

type Rpc struct {
	ChainId uint64 `json:"chainId"`
	Url     string `json:"url"`
}

type NonceManagerConfig struct {
	AutoSync      bool        `json:"autoSync"`
	SyncInterval  uint64      `json:"syncInterval"`
	SyncAfter     uint64      `json:"syncAfter"`
	StorageType   string      `json:"storageType"`
	StorageConfig interface{} `json:"storageConfig"`
	AuthBySig     bool        `json:"authBySig"`
	ApiKey        string      `json:"apiKey"`
}

type WalletManagerConfig struct {
	StorageType   string      `json:"storageType"`
	StorageConfig interface{} `json:"storageConfig"`
}

type Config struct {
	Server        *Server              `json:"server"`
	RpcUrls       []Rpc                `json:"rpcUrls"`
	NonceManager  *NonceManagerConfig  `json:"nonceManager"`
	WalletManager *WalletManagerConfig `json:"walletManager"`
}

func NewConfig(filename string) (*Config, error) {
	config := new(Config)

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	if len(config.Server.LogLevel) == 0 {
		config.Server.LogLevel = "info"
	}

	return config, nil
}
