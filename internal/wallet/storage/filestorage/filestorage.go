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
	types.IWalletStorage
	config *Config
}

func NewFileStorage(c interface{}) (types.IWalletStorage, error) {
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

func (fs *FileStorage) Store(w *types.WalletSerializable) error {
	if len(w.PublicKey) == 0 || len(w.ApiKeyHashed) == 0 || len(w.PrivateKey) == 0 {
		return fmt.Errorf("Wallet object not initiliezed properly")
	}

	data, err := json.Marshal(w)
	if err != nil {
		return err
	}

	filename := fs.getFilename(w.PublicKey)
	logrus.Debugf("Storing wallet info into %s", filename)

	return ioutil.WriteFile(filename, data, 0600)
}

func (fs *FileStorage) Load(address string, apiKeyHashed string) (*types.WalletSerializable, error) {
	filename := fs.getFilename(address)
	logrus.Debugf("Loading wallet info from %s", filename)

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	ws := new(types.WalletSerializable)
	err = json.Unmarshal(data, ws)
	if err != nil {
		return nil, err
	}

	if ws.ApiKeyHashed != apiKeyHashed {
		return nil, fmt.Errorf("Failed to load the wallet")
	}

	return ws, nil
}

func (fs *FileStorage) getFilename(address string) string {
	return path.Join(fs.config.Path, fmt.Sprintf(".%s.wallet.json", address))
}
