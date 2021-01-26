package repository

import (
	"aioz.io/go-aioz/x_gob_explorer/context"
	"aioz.io/go-aioz/x_gob_explorer/domain/db"
	"aioz.io/go-aioz/x_gob_explorer/domain/entity"
	"gorm.io/gorm"
)

type BlockRepo struct {
	db db.Database
}

func NewBlockRepo(db db.Database) BlockRepo {
	return BlockRepo{
		db: db,
	}
}

func (br *BlockRepo) CreateBlock(ctx context.Context, block *entity.Block) error {
	return br.db.PrepareTransaction(
		func(tx *gorm.DB) error {
			return tx.Create(block).Error
		})
}

func (br *BlockRepo) GetPreviousBlockInfo() (*entity.Block, error) {
	b := new(entity.Block)
	err := br.db.GetGormClient().
		Session(&gorm.Session{PrepareStmt: true}).
		Model(&entity.Block{}).
		Order("height DESC").
		First(&b).Error
	if err != nil {
		return nil, err
	}
	return b, nil
}
