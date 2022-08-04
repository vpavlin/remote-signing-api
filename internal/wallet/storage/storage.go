package storage

import (
	"github.com/vpavlin/remote-signing-api/internal/types"
	"github.com/vpavlin/remote-signing-api/internal/wallet/storage/encryptedfilestorage"
	"github.com/vpavlin/remote-signing-api/internal/wallet/storage/filestorage"
)

func NewStorage(storageType string, c interface{}) (types.IWalletStorage, error) {
	var storage types.IWalletStorage
	var err error
	switch storageType {
	case "filestorage":
		storage, err = filestorage.NewFileStorage(c)
		if err != nil {
			return nil, err
		}
	case "encryptedfilestorage":
		storage, err = encryptedfilestorage.NewFileStorage(c)
		if err != nil {
			return nil, err
		}
	}

	return storage, nil
}
