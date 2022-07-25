package wallet

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/vpavlin/remote-signing-api/config"
	"github.com/vpavlin/remote-signing-api/internal/types"
	"github.com/vpavlin/remote-signing-api/internal/wallet/storage"
)

type WalletManager struct {
	wallets []*Wallet
	client  *ethclient.Client
	config  *config.WalletManagerConfig
	storage *types.IWalletStorage
}

func NewWalletManager(client *ethclient.Client, config *config.WalletManagerConfig) (*WalletManager, error) {
	wm := new(WalletManager)
	wm.wallets = make([]*Wallet, 0)
	wm.client = client
	wm.config = config

	storage, err := storage.NewStorage(wm.config.StorageType, wm.config.StorageConfig)
	if err != nil {
		return nil, err
	}

	wm.storage = &storage
	return wm, nil
}

func (wm *WalletManager) New(apiKey string) (*Wallet, error) {
	wallet, err := NewWalletFromKey(wm.storage, apiKey, nil)
	if err != nil {
		return nil, fmt.Errorf("WalletManager: %s", err)
	}

	wm.wallets = append(wm.wallets, wallet)
	return wallet, nil
}

func (wm *WalletManager) GetByAddress(apiKey string, address common.Address) (*Wallet, error) {
	w, err := NewWalletFromStorage(wm.storage, apiKey, address)
	if err != nil {
		return nil, err
	}

	return w, nil
}
