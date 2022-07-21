package nonce

import (
	"log"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ChainID uint64
type Address string

type AddressedNonces struct {
	Nonces map[Address]*Nonce
}

type NonceManager struct {
	ChainNonces map[ChainID]*AddressedNonces
	client      *ethclient.Client
	lock        sync.Mutex
}

func NewNonceManager(client *ethclient.Client) (*NonceManager, error) {
	nm := new(NonceManager)
	nm.client = client
	nm.ChainNonces = make(map[ChainID]*AddressedNonces, 0)

	return nm, nil
}

func (nm *NonceManager) GetNonce(chainId ChainID, address Address) (uint64, error) {
	nm.lock.Lock()
	defer nm.lock.Unlock()
	nonce, err := nm.getNonceObject(chainId, address)
	if err != nil {
		return 0, err
	}

	return nonce.GetNonce()
}

func (nm *NonceManager) ReturnNonce(returnedNonce uint64, chainId ChainID, address Address) error {
	nonce, err := nm.getNonceObject(chainId, address)
	if err != nil {
		return err
	}

	return nonce.ReturnNonce(returnedNonce)
}

func (nm *NonceManager) DecreaseNonce(chainId ChainID, address Address) error {
	nonce, err := nm.getNonceObject(chainId, address)
	if err != nil {
		return err
	}

	return nonce.DecreaseNonce()
}

func (nm *NonceManager) Sync(chainId ChainID, address Address) error {
	nonce, err := nm.getNonceObject(chainId, address)
	if err != nil {
		return err
	}

	return nonce.Sync(nm.client)
}

func (nm *NonceManager) getNonceObject(chainId ChainID, address Address) (*Nonce, error) {
	an, ok := nm.ChainNonces[chainId]
	if !ok {
		log.Println("Could not find chainID")
		an = new(AddressedNonces)
		an.Nonces = make(map[Address]*Nonce)
		nm.ChainNonces[chainId] = an
	}

	nonce, ok := an.Nonces[address]
	if !ok {
		log.Println("Could not find address")

		var err error
		nonce, err = NewNonce(nm.client, common.HexToAddress(string(address)), uint64(chainId))
		if err != nil {
			return nil, err
		}
		an.Nonces[address] = nonce
	}

	return nonce, nil
}
