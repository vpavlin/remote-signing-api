package nonce

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	logrus "github.com/sirupsen/logrus"
	"github.com/vpavlin/remote-signing-api/internal/types"
)

const (
	auto_sync_interval = 10 * time.Minute
	auto_sync_after    = 1 * time.Hour
)

type INonce interface {
	GetNonce() (uint64, error)
	Sync(client *ethclient.Client) error
	ReturnNonce(nonce uint64) error
}

type Nonce struct {
	INonce
	Address        string
	ChainId        uint64
	nonce          uint64
	returnedNonces types.SortedNonceArr
	lock           sync.Mutex
	lastUsed       int64
	storage        types.INonceStorage
}

func NewNonce(client *ethclient.Client, storage *types.INonceStorage, address common.Address, chainId uint64, autoSync bool) (*Nonce, error) {
	return NewNonceWithConfig(client, storage, address, chainId, autoSync, auto_sync_interval, auto_sync_after)
}

func NewNonceWithConfig(client *ethclient.Client, storage *types.INonceStorage, address common.Address, chainId uint64, autoSync bool, syncInterval time.Duration, syncAfter time.Duration) (*Nonce, error) {
	nonce := new(Nonce)
	nonce.Address = address.String()
	nonce.ChainId = chainId
	nonce.storage = *storage

	defer func() {
		if autoSync {
			logrus.WithFields(logrus.Fields{
				"address": address,
				"chainId": chainId,
			}).Info("Starting auto sync")
			go nonce.autoSync(client, syncInterval, syncAfter)
		}
	}()

	ns, err := nonce.storage.Load(chainId, address.String())
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
	} else {
		nonce.deserialize(ns)
		return nonce, nil
	}

	err = nonce.Sync(client)
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

func (n *Nonce) Sync(client *ethclient.Client) (err error) {
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

	nonce, err := client.PendingNonceAt(context.Background(), common.HexToAddress(n.Address))
	if err != nil {
		return fmt.Errorf("GetNonce: cannot get account %s nonce, err: %s, set it to nil!",
			n.Address, err)
	}
	n.nonce = nonce
	n.returnedNonces = make(types.SortedNonceArr, 0)

	return err
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

func (n *Nonce) autoSync(client *ethclient.Client, syncInterval time.Duration, syncAfter time.Duration) error {
	logrus.Info("Starting AutoSync")
	for {
		select {
		case <-time.After(syncInterval):
			//logrus.Debug(n.lastUsed, " + ", int64(syncAfter.Seconds()), " = ", n.lastUsed+int64(syncAfter.Seconds()), " < ", time.Now().Unix())
			if n.lastUsed+int64(syncAfter.Seconds()) < time.Now().Unix() {
				logrus.WithFields(logrus.Fields{
					"address": n.Address,
					"chainId": n.ChainId,
				}).Info("Executing auto-sync")
				err := n.Sync(client)
				if err != nil {
					logrus.Errorf("AutoSync: %s\n", err)
					//return err //TODO: should not return
				}
			}
		}
	}
}

func (n *Nonce) serialize() *types.NonceSerializable {
	ns := new(types.NonceSerializable)

	ns.Address = n.Address
	ns.ChainId = n.ChainId
	ns.Nonce = n.nonce
	ns.LastUsed = n.lastUsed
	ns.ReturnedNonces = n.returnedNonces

	return ns
}

func (n *Nonce) deserialize(ns *types.NonceSerializable) {
	n.Address = ns.Address
	n.ChainId = ns.ChainId
	n.nonce = ns.Nonce
	n.lastUsed = ns.LastUsed
	n.returnedNonces = ns.ReturnedNonces
}

func (n *Nonce) updateLastUsed() {
	n.lastUsed = time.Now().Unix()
}
