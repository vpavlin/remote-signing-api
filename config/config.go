package config

import (
	"encoding/json"
	"io/ioutil"
)

type Server struct {
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
	LogLevel string `json:"logLevel"`
}

type Rpc struct {
	ChainId uint64 `json:"chainId"`
	Url     string `json:"url"`
}

type NonceManagerConfig struct {
	AutoSync     bool   `json:"autoSync"`
	SyncInterval uint64 `json:"syncInterval"`
	SyncAfter    uint64 `json:"syncAfter"`
}

type Config struct {
	Server       *Server             `json:"server"`
	RpcUrls      []Rpc               `json:"rpcUrls"`
	NonceManager *NonceManagerConfig `json:"nonceManager"`
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
