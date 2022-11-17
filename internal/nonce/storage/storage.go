package storage

import (
	"github.com/sirupsen/logrus"
	"github.com/vpavlin/remote-signing-api/internal/nonce/storage/filestorage"
	"github.com/vpavlin/remote-signing-api/internal/nonce/storage/postgres"
	"github.com/vpavlin/remote-signing-api/internal/types"
)

func NewStorage(storageType string, c interface{}) (types.INonceStorage, error) {
	logrus.Infof("Storage type: %s, config: %v", storageType, c)
	var storage types.INonceStorage
	var err error
	switch storageType {
	case "filestorage":
		storage, err = filestorage.NewFileStorage(c)
	case "postgres":
		storage, err = postgres.NewPostgresRepository(c)

	}

	return storage, err
}
