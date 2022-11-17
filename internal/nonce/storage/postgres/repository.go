package postgres

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"os"

	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/vpavlin/remote-signing-api/internal/types"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Nonce struct {
	ID             string `gorm:"primaryKey"`
	Address        string
	Contract       sql.NullString `gorm:"default:NULL"`
	ChainId        uint64
	Nonce          uint64
	ReturnedNonces pq.Int64Array `gorm:"type:integer[]"`
	LastUsed       int64
	SigNonce       bool
}

type PostgresRepository struct {
	types.INonceStorage
	config *Config
	db     *gorm.DB
}

func NewPostgresRepository(c interface{}) (types.INonceStorage, error) {
	config, _ := NewConfig(c)

	db, err := gorm.Open(postgres.Open(config.Connection), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	logrus.Infof("Migrating...")
	err = db.AutoMigrate(&Nonce{})
	if err != nil {
		return nil, err
	}

	return &PostgresRepository{config: config, db: db}, nil
}

func (r *PostgresRepository) Store(n *types.NonceSerializable) error {
	if len(n.Address) == 0 || n.ChainId == 0 {
		return fmt.Errorf("Nonce object not initiliezed properly")
	}

	contract := sql.NullString{}

	if n.Contract != nil {
		contract.String = *n.Contract
		contract.Valid = true
	}

	model := Nonce{
		ID:             getId(n.Address, n.Contract, n.ChainId),
		Address:        n.Address,
		Contract:       contract,
		ChainId:        n.ChainId,
		Nonce:          n.Nonce,
		ReturnedNonces: make(pq.Int64Array, len(n.ReturnedNonces)),
		LastUsed:       n.LastUsed,
		SigNonce:       n.Contract != nil,
	}

	for i, v := range n.ReturnedNonces {
		model.ReturnedNonces[i] = int64(v)
	}

	result := r.db.Save(&model)

	logrus.Infof("Saved: %d", result.RowsAffected)

	return result.Error
}

func (r *PostgresRepository) Load(chainId uint64, address string, contract *string) (*types.NonceSerializable, error) {
	m := Nonce{}

	result := r.db.Find(&m, "id = ?", getId(address, contract, chainId))
	if result.Error != nil {
		logrus.Errorf("Error %s", result.Error)
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, os.ErrNotExist
	}

	logrus.Infof("Contract: %v", m.Contract)
	ns := &types.NonceSerializable{
		Address:        m.Address,
		ChainId:        m.ChainId,
		Nonce:          m.Nonce,
		ReturnedNonces: make(types.SortedNonceArr, len(m.ReturnedNonces)),
		LastUsed:       m.LastUsed,
	}

	if m.Contract.Valid {
		ns.Contract = &m.Contract.String
	}

	for i, v := range m.ReturnedNonces {
		ns.ReturnedNonces[i] = uint64(v)
	}

	return ns, nil
}

func getId(address string, contract *string, chainId uint64) string {
	contractVal := ""
	if contract != nil {
		contractVal = *contract
	}
	hash := sha256.Sum256([]byte(fmt.Sprintf("%s-%s-%d", address, contractVal, chainId)))
	return fmt.Sprintf("%x", hash)
}
