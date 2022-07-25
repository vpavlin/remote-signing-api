package wallet

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	iTypes "github.com/vpavlin/remote-signing-api/internal/types"
)

type Wallet struct {
	ApiKey    string
	PublicKey common.Address

	privateKey *ecdsa.PrivateKey

	storage iTypes.IWalletStorage
}

func NewKey() (*ecdsa.PrivateKey, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func NewWalletFromStorage(storage *iTypes.IWalletStorage, apiKey string, address common.Address) (*Wallet, error) {
	wallet := new(Wallet)
	wallet.storage = *storage
	wallet.PublicKey = address

	apiKeyHashed := crypto.Keccak256Hash([]byte(apiKey))

	err := wallet.load(apiKeyHashed.String())
	if err != nil {
		return nil, err
	}

	if wallet.ApiKey != apiKeyHashed.String() {
		return nil, fmt.Errorf("API Key does not match %s != %s", wallet.ApiKey, apiKeyHashed.String())
	}

	return wallet, nil
}

func NewWalletFromString(storage *iTypes.IWalletStorage, apiKey string, key string) (*Wallet, error) {
	privateKey, err := crypto.HexToECDSA(key)
	if err != nil {
		return nil, err
	}

	return NewWalletFromKey(storage, apiKey, privateKey)
}

func NewWalletFromKey(storage *iTypes.IWalletStorage, apiKey string, key *ecdsa.PrivateKey) (*Wallet, error) {
	if len(apiKey) < 32 {
		return nil, fmt.Errorf("API key length must be at least 32")
	}
	if key == nil {
		var err error
		key, err = NewKey()
		if err != nil {
			return nil, err
		}
	}

	publicKey := key.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("Failed to create new wallet")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	wallet := new(Wallet)

	wallet.storage = *storage

	apiKeyHashed := crypto.Keccak256Hash([]byte(apiKey))
	wallet.ApiKey = apiKeyHashed.String()
	wallet.privateKey = key
	wallet.PublicKey = fromAddress
	wallet.store()

	return wallet, nil
}

func (w *Wallet) GetTxOpts(chainId *big.Int) (*bind.TransactOpts, error) {
	keyAddr := crypto.PubkeyToAddress(w.privateKey.PublicKey)
	if chainId == nil {
		return nil, bind.ErrNoChainID
	}

	return &bind.TransactOpts{
		From: keyAddr,
		Signer: func(address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			return w.SignTX(tx, chainId)
		},
		Context: context.Background(),
	}, nil
}

func (w *Wallet) Sign(data []byte) ([]byte, error) {
	return crypto.Sign(data, w.privateKey)
}

func (w *Wallet) SignTX(tx *types.Transaction, chainId *big.Int) (*types.Transaction, error) {
	signer := types.LatestSignerForChainID(chainId)
	signature, err := w.Sign(signer.Hash(tx).Bytes())

	if err != nil {
		return nil, err
	}

	signedTx, err := tx.WithSignature(signer, signature)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}

func (w *Wallet) ReplaceKey(newApiKey string) error {
	apiKeyHashed := crypto.Keccak256Hash([]byte(newApiKey))

	w.ApiKey = apiKeyHashed.String()

	return w.store()
}

func (w *Wallet) ToSerializable() *iTypes.WalletSerializable {
	return &iTypes.WalletSerializable{
		ApiKeyHashed: w.ApiKey,
		PublicKey:    w.PublicKey.String(),
		PrivateKey:   hex.EncodeToString(crypto.FromECDSA(w.privateKey)),
	}
}

func (w *Wallet) FromSerializable(ws *iTypes.WalletSerializable) error {
	w.ApiKey = ws.ApiKeyHashed
	w.PublicKey = common.HexToAddress(ws.PublicKey)

	privateKey, err := crypto.HexToECDSA(ws.PrivateKey)
	if err != nil {
		return err
	}

	w.privateKey = privateKey

	return nil
}

func (w *Wallet) store() error {
	ws := w.ToSerializable()
	return w.storage.Store(ws)
}

func (w *Wallet) load(apiKeyHashed string) error {
	ws, err := w.storage.Load(w.PublicKey.String(), apiKeyHashed)
	if err != nil {
		return err
	}

	err = w.FromSerializable(ws)
	if err != nil {
		return fmt.Errorf("Wallet: Could not deserialize wallet %s: %s", w.PublicKey, err)
	}

	return nil
}
