package repository

import (
	"aioz.io/go-aioz/x_gob_explorer/context"
	"aioz.io/go-aioz/x_gob_explorer/domain/db"
	"aioz.io/go-aioz/x_gob_explorer/domain/entity"
	"gorm.io/gorm"
)

type DelegatorRepo struct {
	db db.Database
}

func NewDelegatorRepo(db db.Database) DelegatorRepo {
	return DelegatorRepo{
		db: db,
	}
}

func (dr *DelegatorRepo) CreateDelegator(ctx context.Context, dele *entity.Delegator) error {
	return dr.db.PrepareTransaction(
		func(tx *gorm.DB) error {
			return tx.Session(&gorm.Session{PrepareStmt: true}).Create(dele).Error
		})
}

func (dr *DelegatorRepo) UpsertDelegator(ctx context.Context, obj *entity.Delegator) error {
	return dr.db.PrepareTransaction(
		func(tx *gorm.DB) error {
			return tx.Session(&gorm.Session{PrepareStmt: true}).Exec(`UPSERT INTO 
			delegators(delegator_address,validator_address,shares) 
			VALUES(?,?,?)`,
				obj.DelegatorAddress, obj.ValidatorAddress, obj.Shares).Error
		})
}

func (dr *DelegatorRepo) UpdateDelegatorWithdrawAddress(ctx context.Context, obj *entity.Delegator) error {
	return dr.db.PrepareTransaction(
		func(tx *gorm.DB) error {
			return tx.Session(&gorm.Session{PrepareStmt: true}).Model(&entity.Delegator{}).
				Where("delegator_address = ?", obj.DelegatorAddress).
				Update("withdraw_address", obj.WithdrawAddress).
				Error
		})
}

func (dr *DelegatorRepo) GetDelegatorByDelegatorAndValidator(ctx context.Context, delegatorAddress, validatorAddress string) (*entity.Delegator, error) {
	dele := new(entity.Delegator)
	err := dr.db.GetGormClient().
		Session(&gorm.Session{PrepareStmt: true}).
		Model(&entity.Delegator{}).
		Where("delegator_address = ? AND validator_address = ?", delegatorAddress, validatorAddress).
		First(&dele).
		Error
	return dele, err
}

func (dr *DelegatorRepo) GetDelegatorByValidator(ctx context.Context, validatorAddress string) ([]*entity.Delegator, error) {
	dele := make([]*entity.Delegator, 0)
	err := dr.db.GetGormClient().
		Session(&gorm.Session{PrepareStmt: true}).
		Model(&entity.Delegator{}).
		Where("validator_address = ?", validatorAddress).
		Find(&dele).
		Error
	return dele, err
}
