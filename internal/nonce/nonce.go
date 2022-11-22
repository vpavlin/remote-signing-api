package nonce

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	logrus "github.com/sirupsen/logrus"
	"github.com/vpavlin/remote-signing-api/internal/bindings"
	"github.com/vpavlin/remote-signing-api/internal/types"
)

const (
	auto_sync_max_interval = 24 * time.Hour
	auto_sync_interval     = 10 * time.Minute
	auto_sync_after        = 1 * time.Hour
	auto_sync_max_err      = 10
)

type INonce interface {
	GetNonce() (uint64, error)
	Sync(client *ethclient.Client) error
	ReturnNonce(nonce uint64) error
}

type Nonce struct {
	INonce
	Address        string
	Contract       *string
	ChainId        uint64
	nonce          uint64
	returnedNonces types.SortedNonceArr
	lock           sync.Mutex
	lastUsed       int64
	storage        types.INonceStorage
	ErrCount       uint
}

func NewNonce(client *ethclient.Client, storage *types.INonceStorage, address string, chainId uint64, contract *string, autoSync bool) (*Nonce, error) {
	return NewNonceWithConfig(client, storage, address, chainId, contract, autoSync, auto_sync_interval, auto_sync_after)
}

func NewNonceWithConfig(client *ethclient.Client, storage *types.INonceStorage, address string, chainId uint64, contract *string, autoSync bool, syncInterval time.Duration, syncAfter time.Duration) (nonce *Nonce, err error) {

	nonce = &Nonce{
		Address:  address,
		ChainId:  chainId,
		storage:  *storage,
		Contract: contract,
	}

	defer func() {
		logrus.Infof("Error: %s", err)
		if err == nil && autoSync {
			logrus.WithFields(logrus.Fields{
				"address": address,
				"chainId": chainId,
			}).Info("Starting auto sync")
			go nonce.autoSync(client, syncInterval, syncAfter)
		}
	}()

	err = nonce.Load()
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
	} else {
		return nonce, nil
	}

	_, err = nonce.Sync(client)
	if err != nil {
		return nil, err
	}
	nonce.returnedNonces = make(types.SortedNonceArr, 0)

	err = nonce.storage.Store(nonce.serialize())
	if err != nil {
		return nil, err
	}

	return nonce, nil
}

func (n *Nonce) Sync(client *ethclient.Client) (updated bool, err error) {
	n.lock.Lock()
	defer n.lock.Unlock()

	defer func() {
		if updated {
			deferErr := n.storage.Store(n.serialize())
			if deferErr != nil {
				logrus.Error(deferErr)
				err = deferErr
			}
		}
	}()

	n.updateLastUsed()
	var nonce uint64

	if n.Contract != nil {
		c, err := bindings.NewSigNonce(common.HexToAddress(*n.Contract), client)
		if err != nil {
			return false, fmt.Errorf("Failed to bound the contract %s: %s", *n.Contract, err)
		}
		bigNonce, err := c.SigNonce(&bind.CallOpts{}, common.HexToAddress(n.Address))
		if err != nil {
			return false, fmt.Errorf("Failed to get nonce from %s: %s", *n.Contract, err)
		}

		nonce = bigNonce.Uint64()
	} else {
		nonce, err = client.PendingNonceAt(context.Background(), common.HexToAddress(n.Address))
		if err != nil {
			return false, fmt.Errorf("GetNonce: cannot get account %s nonce, err: %s, set it to nil!",
				n.Address, err)
		}
	}

	if nonce == n.nonce {
		return false, nil
	}
	updated = true
	n.nonce = nonce
	n.returnedNonces = make(types.SortedNonceArr, 0)

	return true, err
}

func (n *Nonce) GetNonce() (nonce uint64, err error) {
	n.lock.Lock()
	defer n.lock.Unlock()
	defer func() {
		deferErr := n.storage.Store(n.serialize())
		if deferErr != nil {
			logrus.Error(deferErr)
			err = deferErr
		}
	}()

	n.updateLastUsed()

	if n.returnedNonces.Len() > 0 {
		logrus.Debug("Using returned nonces: ", n.returnedNonces)
		nonce := n.returnedNonces[0]
		n.returnedNonces = n.returnedNonces[1:]
		return nonce, nil
	}

	// return a new point
	nonce = n.nonce
	logrus.Debug("Using nonce: ", n.nonce)

	// increase record
	n.nonce++

	return nonce, nil
}

func (n *Nonce) ReturnNonce(nonce uint64) (err error) {
	n.lock.Lock()
	defer n.lock.Unlock()
	defer func() {
		deferErr := n.storage.Store(n.serialize())
		if deferErr != nil {
			logrus.Error(deferErr)
			err = deferErr
		}
	}()
	n.updateLastUsed()

	if nonce >= n.nonce {
		return fmt.Errorf("Returned nonce too high")
	}

	for _, returned := range n.returnedNonces {
		if returned == nonce {
			return fmt.Errorf("Nonce already returned")
		}
	}
	arr := n.returnedNonces
	arr = append(arr, nonce)
	sort.Sort(arr)
	n.returnedNonces = arr

	return nil

}

func (n *Nonce) DecreaseNonce() error {
	var err error
	n.lock.Lock()
	defer n.lock.Unlock()
	defer func() {
		deferErr := n.storage.Store(n.serialize())
		if deferErr != nil {
			err = deferErr
		}
	}()
	n.updateLastUsed()

	if n.nonce > 0 {
		n.nonce--
	}

	return err

}

func (n *Nonce) Load() error {
	ns, err := n.storage.Load(n.ChainId, n.Address, n.Contract)
	if err != nil {
		return err
	}

	n.deserialize(ns)
	return nil
}

func (n *Nonce) autoSync(client *ethclient.Client, syncInterval time.Duration, syncAfter time.Duration) error {
	logger := logrus.WithFields(logrus.Fields{
		"address": n.Address,
		"chainId": n.ChainId,
	})
	if n.Contract != nil {
		logger = logger.WithField("contract", *n.Contract)
	}

	interval := syncInterval
	logger.Info("Starting AutoSync")
	for {
		select {
		case <-time.After(interval):
			//logrus.Debug(n.lastUsed, " + ", int64(syncAfter.Seconds()), " = ", n.lastUsed+int64(syncAfter.Seconds()), " < ", time.Now().Unix())
			if n.lastUsed+int64(syncAfter.Seconds()) < time.Now().Unix() {
				logger.Infof("Executing auto-sync (interval: %d)", interval)
				updated, err := n.Sync(client)
				if err != nil {
					interval *= 2
					logger.Errorf("AutoSync: %s\n", err)
					n.ErrCount++
					if n.ErrCount > auto_sync_max_err {
						return fmt.Errorf("Stopping the auto sync: %s", err)
					}
					//return err //TODO: should not return
				}

				if updated {
					interval = syncInterval
				} else if interval < auto_sync_max_interval {
					interval *= 2
				}
			}
		}
	}
}

func (n *Nonce) serialize() *types.NonceSerializable {
	ns := new(types.NonceSerializable)

	ns.Address = n.Address
	ns.Contract = n.Contract
	ns.ChainId = n.ChainId
	ns.Nonce = n.nonce
	ns.LastUsed = n.lastUsed
	ns.ReturnedNonces = n.returnedNonces

	return ns
}

func (n *Nonce) deserialize(ns *types.NonceSerializable) {
	n.Address = ns.Address
	n.Contract = ns.Contract
	n.ChainId = ns.ChainId
	n.nonce = ns.Nonce
	n.lastUsed = ns.LastUsed
	n.returnedNonces = ns.ReturnedNonces
}

func (n *Nonce) updateLastUsed() {
	n.lastUsed = time.Now().Unix()
}
