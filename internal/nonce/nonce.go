package nonce

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type INonce interface {
	GetNonce() (uint64, error)
	DecreaseNonce()
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
}

func NewNonce(client *ethclient.Client, address common.Address, chainId uint64) (*Nonce, error) {
	nonce := new(Nonce)
	nonce.Address = address.String()
	nonce.ChainId = chainId

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
		log.Println("Using returned nonces: ", n.returnedNonces)
		nonce := n.returnedNonces[0]
		n.returnedNonces = n.returnedNonces[1:]
		log.Println(n.returnedNonces)
		return nonce, nil
	}

	// return a new point
	nonce := n.nonce
	log.Println("Using nonce: ", n.nonce)

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
}

func (n *Nonce) store() error {
	if len(n.Address) == 0 || n.ChainId == 0 {
		return fmt.Errorf("Nonce object not initiliezed properly")
	}

	ns := new(NonceSerializable)
	ns.Address = n.Address
	ns.ChainId = n.ChainId
	ns.Nonce = n.nonce
	ns.ReturnedNonces = n.returnedNonces

	data, err := json.Marshal(ns)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf(".%s-%d.nonce.json", n.Address, n.ChainId)
	return ioutil.WriteFile(filename, data, 0600)
}

func (n *Nonce) load() error {
	filename := fmt.Sprintf(".%s-%d.nonce.json", n.Address, n.ChainId)

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

	return nil
}

type SortedNonceArr []uint64

func (arr SortedNonceArr) Less(i, j int) bool {
	return arr[i] < arr[j]
}

func (arr SortedNonceArr) Len() int { return len(arr) }

func (arr SortedNonceArr) Swap(i, j int) { arr[i], arr[j] = arr[j], arr[i] }
