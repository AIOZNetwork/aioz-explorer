package repository

import (
	"aioz.io/go-aioz/x_gob_explorer/context"
	"aioz.io/go-aioz/x_gob_explorer/domain/db"
	"aioz.io/go-aioz/x_gob_explorer/domain/entity"
	"fmt"
	"gorm.io/gorm"
	"log"
	"time"
)

type WalletRepo struct {
	db db.Database
}

func NewWalletRepo(db db.Database) WalletRepo {
	return WalletRepo{
		db: db,
	}
}

func (wr *WalletRepo) UpsertWallet(ctx context.Context, wl *entity.WalletAddress) error {
	return wr.db.PrepareTransaction(
		func(tx *gorm.DB) error {
			return tx.Session(&gorm.Session{PrepareStmt: true}).Exec(`UPSERT INTO 
			wallet_addresses(address,coins,account_number,sequence,pub_key,coin_aioz,coin_stake) 
			VALUES(?,?,?,?,?,?,?)`,
				wl.Address, wl.Coins, wl.AccountNumber, wl.Sequence, wl.PubKey, wl.CoinAIOZ, wl.CoinStake).Error
		})
}

func (wr *WalletRepo) MultiRowsUpsert(ctx context.Context, accAddrs []*entity.WalletAddress) error {
	if len(accAddrs) > 0 {
		sql := `UPSERT INTO wallet_addresses(address,coins,account_number,sequence,pub_key,coin_aioz,coin_stake) VALUES`
		for i, t := range accAddrs {
			sql += fmt.Sprintf(`('%v','%v',%v,%v,'%v',%v,%v)`,
				t.Address, t.Coins, t.AccountNumber, t.Sequence, t.PubKey, t.CoinAIOZ, t.CoinStake)
			if i == len(accAddrs)-1 {
				sql += ";"
			} else {
				sql += ","
			}
		}

		return wr.db.GetGormClient().Transaction(
			func(tx *gorm.DB) error {
				if err := tx.Exec(sql).Error; err != nil {
					log.Println(err)
				}
				return nil
			})
	}
	return nil
}

func (wr *WalletRepo) GetVestingWallet(ctx context.Context, currTime time.Time) ([]*entity.WalletAddress, error) {
	result := make([]*entity.WalletAddress, 0)
	err := wr.db.GetGormClient().
		Session(&gorm.Session{PrepareStmt: true}).
		Model(&entity.WalletAddress{}).
		Where("end_time > ?", currTime.Unix()).
		Find(&result).Error
	return result, err
}
