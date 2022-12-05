package encryptedfilestorage

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/sirupsen/logrus"
	"github.com/vpavlin/remote-signing-api/internal/types"
)

const walletFileName = "wallet.dat"

type FileStorage struct {
	types.IWalletStorage
	config   *Config
	password string
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

	password := os.Getenv(config.PasswordEnvName)
	fs := &FileStorage{config: config, password: password}

	if len(password) < 10 {
		return nil, fmt.Errorf("Missing or short password.")
	}

	return fs, nil
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

	encrypted, err := encrypt(data, fs.getFilePassword(w.ApiKeyHashed))
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, encrypted, 0600)
}

func (fs *FileStorage) Load(address string, apiKeyHashed string) (*types.WalletSerializable, error) {
	filename := fs.getFilename(address)
	logrus.Debugf("Loading wallet info from %s", filename)

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	decrypted, err := decrypt(data, fs.getFilePassword(apiKeyHashed))
	if err != nil {
		return nil, err
	}

	ws := new(types.WalletSerializable)
	err = json.Unmarshal(decrypted, ws)
	if err != nil {
		return nil, err
	}

	if ws.ApiKeyHashed != apiKeyHashed {
		return nil, fmt.Errorf("Failed to load the wallet")
	}

	return ws, nil
}

func (fs *FileStorage) getFilename(address string) string {
	return path.Join(fs.config.Path, fmt.Sprintf("%s.wallet.dat", address))
}

func (fs *FileStorage) getFilePassword(key string) []byte {
	keyb := sha256.Sum256([]byte(fmt.Sprintf("%s#%s", key, fs.password)))
	return keyb[:]
}
