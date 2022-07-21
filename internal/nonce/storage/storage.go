package storage

import (
	"github.com/vpavlin/remote-signing-api/internal/nonce/storage/filestorage"
	"github.com/vpavlin/remote-signing-api/internal/types"
)

func NewStorage(storageType string, c interface{}) (types.INonceStorage, error) {
	var storage types.INonceStorage
	var err error
	switch storageType {
	case "filestorage":
		storage, err = filestorage.NewFileStorage(c)
		if err != nil {
			return nil, err
		}
	}

	return storage, nil
}
