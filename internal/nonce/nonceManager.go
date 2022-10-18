package nonce

import (
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
	"github.com/vpavlin/remote-signing-api/config"
	"github.com/vpavlin/remote-signing-api/internal/nonce/storage"
	"github.com/vpavlin/remote-signing-api/internal/types"
)

type ChainID uint64

type AddressedNonces struct {
	Nonces map[string]*Nonce
}

type NonceManager struct {
	ChainNonces map[ChainID]*AddressedNonces
	clients     map[ChainID]*ethclient.Client
	lock        sync.Mutex
	config      *config.NonceManagerConfig
	storage     *types.INonceStorage
}

func NewNonceManager(rpcUrls []config.Rpc, config *config.NonceManagerConfig) (*NonceManager, error) {
	nm := new(NonceManager)
	nm.ChainNonces = make(map[ChainID]*AddressedNonces, 0)
	nm.clients = make(map[ChainID]*ethclient.Client)
	nm.initClients(rpcUrls)
	nm.config = config
	storage, err := storage.NewStorage("filestorage", nm.config.StorageConfig)
	if err != nil {
		return nil, err
	}
	nm.storage = &storage
	return nm, nil
}

func (nm *NonceManager) initClients(rpcUrls []config.Rpc) error {
	for _, rpc := range rpcUrls {
		client, err := ethclient.Dial(rpc.Url)
		if err != nil {
			return err
		}
		nm.clients[ChainID(rpc.ChainId)] = client
	}

	return nil
}

func (nm *NonceManager) GetNonce(chainId ChainID, address string, contract *string) (uint64, error) {
	nonce, err := nm.getNonceObject(chainId, address, contract)
	if err != nil {
		return 0, err
	}

	return nonce.GetNonce()
}

func (nm *NonceManager) ReturnNonce(returnedNonce uint64, chainId ChainID, address string, contract *string) error {
	nonce, err := nm.getNonceObject(chainId, address, contract)
	if err != nil {
		return err
	}

	return nonce.ReturnNonce(returnedNonce)
}

func (nm *NonceManager) DecreaseNonce(chainId ChainID, address string, contract *string) error {
	nonce, err := nm.getNonceObject(chainId, address, contract)
	if err != nil {
		return err
	}

	return nonce.DecreaseNonce()
}

func (nm *NonceManager) Sync(chainId ChainID, address string, contract *string) error {
	nonce, err := nm.getNonceObject(chainId, address, contract)
	if err != nil {
		return err
	}

	client, ok := nm.clients[chainId]
	if !ok {
		return fmt.Errorf("Unknown client for chainId %d", chainId)
	}

	return nonce.Sync(client)
}

func (nm *NonceManager) getNonceObject(chainId ChainID, address string, contract *string) (*Nonce, error) {
	an, ok := nm.ChainNonces[chainId]
	if !ok {
		logrus.WithFields(logrus.Fields{
			"address": address,
			"chainId": chainId,
		}).Info("Initializing new chainID")
		an = new(AddressedNonces)
		an.Nonces = make(map[string]*Nonce)
		nm.ChainNonces[chainId] = an
	}

	nonce, ok := an.Nonces[address]
	if !ok {
		logrus.WithFields(logrus.Fields{
			"address": address,
			"chainId": chainId,
		}).Info("Setting up new nonce")

		client, ok := nm.clients[chainId]
		if !ok {
			return nil, fmt.Errorf("Unknown client for chainId %d", chainId)
		}

		var err error

		syncInterval := time.Duration(nm.config.SyncInterval) * time.Second
		syncAfter := time.Duration(nm.config.SyncAfter) * time.Second

		nonce, err = NewNonceWithConfig(client, nm.storage, address, uint64(chainId), contract, nm.config.AutoSync, syncInterval, syncAfter)
		if err != nil {
			return nil, err
		}
		an.Nonces[address] = nonce
	}

	return nonce, nil
}
