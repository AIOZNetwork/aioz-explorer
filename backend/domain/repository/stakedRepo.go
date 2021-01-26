package repository

import (
	"aioz.io/go-aioz/x_gob_explorer/context"
	"aioz.io/go-aioz/x_gob_explorer/domain/db"
	"aioz.io/go-aioz/x_gob_explorer/domain/entity"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type StakedRepo struct {
	db db.Database
}

func NewStakedRepo(db db.Database) StakedRepo {
	return StakedRepo{
		db: db,
	}
}

func (sr *StakedRepo) GetCurrentStakeByDelegator(ctx context.Context, address string) (*entity.Stake, error) {
	var s *entity.Stake
	rows, err := sr.db.GetGormClient().
		Session(&gorm.Session{PrepareStmt: true}).
		Table("stakes").
		Select("delegator_address, shares, shares_dec").
		Where("delegator_address = ?", address).
		Rows()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var a, shares string
		var sharesDec decimal.Decimal
		if err := rows.Scan(&a, &shares, &sharesDec); err != nil {
			return nil, err
		}
		s = &entity.Stake{
			DelegatorAddress: a,
			Shares:           shares,
			SharesDec:        sharesDec,
		}
		break
	}
	return s, err

	//err := sr.db.GetGormClient().Model(&entity.Stake{}).Where("delegator_address = ?", address).Find(&s).Error
	//return s, err
}

func (sr *StakedRepo) UpsertStakedCoins(ctx context.Context, obj *entity.Stake) error {
	return sr.db.PrepareTransaction(
		func(tx *gorm.DB) error {
			return tx.Session(&gorm.Session{PrepareStmt: true}).Exec(`UPSERT INTO 
			stakes(delegator_address,shares,shares_dec) 
			VALUES(?,?,?)`,
				obj.DelegatorAddress, obj.Shares, obj.SharesDec).Error
		})
}
