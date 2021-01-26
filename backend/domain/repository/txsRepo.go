package repository

import (
	"aioz.io/go-aioz/x_gob_explorer/context"
	"aioz.io/go-aioz/x_gob_explorer/domain/db"
	"aioz.io/go-aioz/x_gob_explorer/domain/entity"
	"fmt"
	"gorm.io/gorm"
	"log"
)

type TxsRepo struct {
	db db.Database
}

func NewTxsRepo(db db.Database) TxsRepo {
	return TxsRepo{
		db: db,
	}
}

func (txr *TxsRepo) CreateTxs(ctx context.Context, txs *entity.Txs) error {
	return txr.db.PrepareTransaction(
		func(tx *gorm.DB) error {
			return tx.Session(&gorm.Session{PrepareStmt: true}).Create(txs).Error
		})
}

func (txr *TxsRepo) MultiRowsInsertTxs(ctx context.Context, txs []*entity.Txs) error {
	if len(txs) > 0 {
		sql := `INSERT INTO txs(block_height,block_hash,block_time,transaction_hash,transaction_index,message_type,message_index,address,address_from,address_to,amount,is_valid,payload_msg,payload_err) VALUES`
		for i, t := range txs {
			addressFrom := `ARRAY[`
			addressTo := `ARRAY[`
			for id, a := range t.AddressFrom {
				if id == len(t.AddressFrom)-1 {
					addressFrom += fmt.Sprintf(`'%v']`, a)
				} else {
					addressFrom += fmt.Sprintf(`'%v',`, a)
				}
			}
			if len(t.AddressFrom) == 0 {
				addressFrom += "]"
			}
			for id, a := range t.AddressTo {
				if id == len(t.AddressTo)-1 {
					addressTo += fmt.Sprintf(`'%v']`, a)
				} else {
					addressTo += fmt.Sprintf(`'%v',`, a)
				}
			}
			if len(t.AddressTo) == 0 {
				addressTo += "]"
			}
			sql += fmt.Sprintf(`(%v,'%v',%v,'%v',%v,'%v',%v,'%v',%v,%v,'%v',%v,'%v','%v')`,
				t.BlockHeight, t.BlockHash, t.BlockTime, t.TransactionHash, t.TransactionIndex, t.MessageType, t.MessageIndex,
				t.Address, addressFrom, addressTo, t.Amount, t.IsValid, t.PayloadMsg, t.PayloadErr)
			if i == len(txs)-1 {
				sql += ";"
			} else {
				sql += ","
			}
		}

		return txr.db.GetGormClient().Transaction(
			func(tx *gorm.DB) error {
				if err := tx.Exec(sql).Error; err != nil {
					log.Println(err)
				}
				return nil
			})
	}
	return nil
}
