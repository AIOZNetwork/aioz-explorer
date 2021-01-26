package repository

import (
	"aioz.io/go-aioz/x_gob_explorer/context"
	"aioz.io/go-aioz/x_gob_explorer/domain/db"
	"aioz.io/go-aioz/x_gob_explorer/domain/entity"
	"fmt"
	"gorm.io/gorm"
)

type ValidatorRepo struct {
	db db.Database
}

func NewValidatorRepo(db db.Database) ValidatorRepo {
	return ValidatorRepo{
		db: db,
	}
}

func (vr *ValidatorRepo) GetAllValidators(ctx context.Context) ([]*entity.Validator, error) {
	result := make([]*entity.Validator, 0)
	err := vr.db.GetGormClient().
		Session(&gorm.Session{PrepareStmt: true}).
		Model(&entity.Validator{}).
		Find(&result).
		Error
	return result, err
}

func (vr *ValidatorRepo) CreateValidator(ctx context.Context, val *entity.Validator) error {
	return vr.db.PrepareTransaction(
		func(tx *gorm.DB) error {
			return tx.Session(&gorm.Session{PrepareStmt: true}).Create(val).Error
		})
}

func (vr *ValidatorRepo) UpsertValidator(ctx context.Context, obj *entity.Validator) error {
	return vr.db.PrepareTransaction(
		func(tx *gorm.DB) error {
			return tx.Session(&gorm.Session{PrepareStmt: true}).Exec(`UPSERT INTO 
			validators(address,tokens,power,jailed,status) 
			VALUES(?,?,?,?,?)`,
				obj.Address, obj.Tokens, obj.Power, obj.Jailed, obj.Status).Error
		})
}

func (vr *ValidatorRepo) MultiRowsUpsert(ctx context.Context, vals []*entity.Validator) error {
	if len(vals) > 0 {
		sql := `UPSERT INTO validators(address,tokens,power,jailed,status) VALUES`
		for i, t := range vals {
			sql += fmt.Sprintf(`('%v','%v',%v,%v,'%v')`,
				t.Address, t.Tokens, t.Power, t.Jailed, t.Status)
			if i == len(vals)-1 {
				sql += ";"
			} else {
				sql += ","
			}
		}

		return vr.db.PrepareTransaction(
			func(tx *gorm.DB) error {
				return tx.Exec(sql).Error
			})
	}
	return nil
}

func (vr *ValidatorRepo) UpdateValidator(ctx context.Context, obj *entity.Validator) error {
	return vr.db.PrepareTransaction(
		func(tx *gorm.DB) error {
			return tx.Model(&entity.Validator{}).
				Where("address = ?", obj.Address).
				Updates(&entity.Validator{
					Website:  obj.Website,
					Moniker:  obj.Moniker,
					Identity: obj.Identity,
					Detail:   obj.Detail,
				}).Error
		})
}

func (vr *ValidatorRepo) UpsertValidatorRewardPool(ctx context.Context, obj *entity.Validator) error {
	return vr.db.PrepareTransaction(
		func(tx *gorm.DB) error {
			return tx.Exec(`UPSERT INTO validators(address,period,reward_pool)
								VALUES(?,?,?)`, obj.Address, obj.Period, obj.RewardPool).Error
		})
}
