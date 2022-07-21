package nonce

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	logrus "github.com/sirupsen/logrus"
)

const (
	auto_sync_interval = 10 * time.Minute
	auto_sync_after    = 1 * time.Hour
)

type INonce interface {
	GetNonce() (uint64, error)
	Sync(client *ethclient.Client) error
	ReturnNonce(nonce uint64) error
	store() error
	load() error
}

type Nonce struct {
	INonce
	Address        string
	ChainId        uint64
	nonce          uint64
	returnedNonces SortedNonceArr
	lock           sync.Mutex
	lastUsed       int64
}

func NewNonce(client *ethclient.Client, address common.Address, chainId uint64, autoSync bool) (*Nonce, error) {
	return NewNonceWithConfig(client, address, chainId, autoSync, auto_sync_interval, auto_sync_after)
}

func NewNonceWithConfig(client *ethclient.Client, address common.Address, chainId uint64, autoSync bool, syncInterval time.Duration, syncAfter time.Duration) (*Nonce, error) {
	nonce := new(Nonce)
	nonce.Address = address.String()
	nonce.ChainId = chainId

	defer func() {
		if autoSync {
			logrus.WithFields(logrus.Fields{
				"address": address,
				"chainId": chainId,
			}).Info("Starting auto sync")
			go nonce.autoSync(client, syncInterval, syncAfter)
		}
	}()

	err := nonce.load()
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
	} else {
		return nonce, nil
	}

	err = nonce.Sync(client)
	if err != nil {
		return nil, err
	}
	nonce.returnedNonces = make(SortedNonceArr, 0)

	nonce.store()

	return nonce, nil
}

func (n *Nonce) Sync(client *ethclient.Client) error {
	n.lock.Lock()
	defer n.lock.Unlock()
	defer n.store()

	nonce, nonceErr := client.PendingNonceAt(context.Background(), common.HexToAddress(n.Address))
	if nonceErr != nil {
		return fmt.Errorf("GetNonce: cannot get account %s nonce, err: %s, set it to nil!",
			n.Address, nonceErr)
	}
	n.nonce = nonce
	n.returnedNonces = make(SortedNonceArr, 0)

	return nil
}

func (n *Nonce) GetNonce() (uint64, error) {
	n.lock.Lock()
	defer n.lock.Unlock()
	defer n.store()

	if n.returnedNonces.Len() > 0 {
		logrus.Debug("Using returned nonces: ", n.returnedNonces)
		nonce := n.returnedNonces[0]
		n.returnedNonces = n.returnedNonces[1:]
		return nonce, nil
	}

	// return a new point
	nonce := n.nonce
	logrus.Debug("Using nonce: ", n.nonce)

	// increase record
	n.nonce++

	return nonce, nil
}

func (n *Nonce) ReturnNonce(nonce uint64) error {
	n.lock.Lock()
	defer n.lock.Unlock()
	defer n.store()

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
		deferErr := n.store()
		if deferErr != nil {
			err = deferErr
		}
	}()

	if n.nonce > 0 {
		n.nonce--
	}

	return err

}

type NonceSerializable struct {
	Address        string         `json:"address"`
	ChainId        uint64         `json:"chainId"`
	Nonce          uint64         `json:"nonce"`
	ReturnedNonces SortedNonceArr `json:"returnedNonces"`
	LastUsed       int64          `json:"lastUsed"`
}

func (n *Nonce) store() error {
	if len(n.Address) == 0 || n.ChainId == 0 {
		return fmt.Errorf("Nonce object not initiliezed properly")
	}

	n.lastUsed = time.Now().Unix()

	ns := new(NonceSerializable)
	ns.Address = n.Address
	ns.ChainId = n.ChainId
	ns.Nonce = n.nonce
	ns.ReturnedNonces = n.returnedNonces
	ns.LastUsed = n.lastUsed

	data, err := json.Marshal(ns)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf(".%s-%d.nonce.json", n.Address, n.ChainId)
	logrus.Debugf("Storing nonce info into %s", filename)

	return ioutil.WriteFile(filename, data, 0600)
}

func (n *Nonce) load() error {
	filename := fmt.Sprintf(".%s-%d.nonce.json", n.Address, n.ChainId)

	logrus.Debugf("Loading nonce info from %s", filename)

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	ns := new(NonceSerializable)
	err = json.Unmarshal(data, ns)
	if err != nil {
		return err
	}

	n.nonce = ns.Nonce
	n.returnedNonces = ns.ReturnedNonces
	n.lastUsed = ns.LastUsed

	return nil
}

func (n *Nonce) autoSync(client *ethclient.Client, syncInterval time.Duration, syncAfter time.Duration) error {
	for {
		select {
		case <-time.After(syncInterval):
			logrus.Debug(n.lastUsed, " + ", int64(syncAfter.Seconds()), " = ", n.lastUsed+int64(syncAfter.Seconds()), " < ", time.Now().Unix())
			if n.lastUsed+int64(syncAfter.Seconds()) < time.Now().Unix() {
				logrus.WithFields(logrus.Fields{
					"address": n.Address,
					"chainId": n.ChainId,
				}).Info("Executing auto-sync")
				n.Sync(client)
			}
			//log.Infof("clearNonce: clear all cache nonce")

		}
	}
}

type SortedNonceArr []uint64

func (arr SortedNonceArr) Less(i, j int) bool {
	return arr[i] < arr[j]
}

func (arr SortedNonceArr) Len() int { return len(arr) }

func (arr SortedNonceArr) Swap(i, j int) { arr[i], arr[j] = arr[j], arr[i] }
