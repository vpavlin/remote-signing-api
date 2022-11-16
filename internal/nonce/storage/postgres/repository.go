package postgres

import (
	"fmt"

	"github.com/lib/pq"
	"github.com/vpavlin/remote-signing-api/internal/types"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Nonce struct {
	Address        string  `gorm:"primaryKey"`
	Contract       *string `gorm:"primaryKey"`
	ChainId        uint64  `gorm:"primaryKey"`
	Nonce          uint64
	ReturnedNonces pq.Int64Array `gorm:"type:integer[]"`
	LastUsed       int64
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

	model := Nonce{
		Address:        n.Address,
		Contract:       n.Contract,
		ChainId:        n.ChainId,
		Nonce:          n.Nonce,
		ReturnedNonces: make(pq.Int64Array, len(n.ReturnedNonces)),
		LastUsed:       n.LastUsed,
	}

	for i, v := range n.ReturnedNonces {
		model.ReturnedNonces[i] = int64(v)
	}

	result := r.db.Save(&model)

	return result.Error
}

func (r *PostgresRepository) Load(chainId uint64, address string, contract *string) (*types.NonceSerializable, error) {
	m := Nonce{}
	result := r.db.Find(&m, "address = ? && contract = ? && chainId = ?", address, contract, chainId)
	if result.Error != nil {
		return nil, result.Error
	}
	ns := &types.NonceSerializable{
		Address:        m.Address,
		Contract:       m.Contract,
		ChainId:        m.ChainId,
		Nonce:          m.Nonce,
		ReturnedNonces: make(types.SortedNonceArr, len(m.ReturnedNonces)),
		LastUsed:       m.LastUsed,
	}

	for i, v := range m.ReturnedNonces {
		ns.ReturnedNonces[i] = uint64(v)
	}

	return ns, nil
}
