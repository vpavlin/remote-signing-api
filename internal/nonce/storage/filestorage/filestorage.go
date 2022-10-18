package filestorage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/sirupsen/logrus"
	"github.com/vpavlin/remote-signing-api/internal/types"
)

type FileStorage struct {
	types.INonceStorage
	config *Config
}

func NewFileStorage(c interface{}) (types.INonceStorage, error) {
	config, _ := NewConfig(c)

	_, err := os.Stat(config.Path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = os.MkdirAll(config.Path, 0700)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return &FileStorage{config: config}, nil
}

func (fs *FileStorage) Store(n *types.NonceSerializable) error {
	if len(n.Address) == 0 || n.ChainId == 0 {
		return fmt.Errorf("Nonce object not initiliezed properly")
	}

	data, err := json.Marshal(n)
	if err != nil {
		return err
	}

	filename := fs.getFilename(n.ChainId, n.Address, n.Contract)
	logrus.Debugf("Storing nonce info into %s", filename)

	return ioutil.WriteFile(filename, data, 0600)
}

func (fs *FileStorage) Load(chainId uint64, address string, contract *string) (*types.NonceSerializable, error) {
	filename := fs.getFilename(chainId, address, contract)
	logrus.Debugf("Loading nonce info from %s", filename)

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	ns := new(types.NonceSerializable)
	err = json.Unmarshal(data, ns)
	if err != nil {
		return nil, err
	}

	return ns, nil
}

func (fs *FileStorage) getFilename(chainId uint64, address string, contract *string) string {
	if contract != nil {
		return path.Join(fs.config.Path, fmt.Sprintf(".%s-%s-%d.nonce.json", address, *contract, chainId))
	}

	return path.Join(fs.config.Path, fmt.Sprintf(".%s-%d.nonce.json", address, chainId))
}
