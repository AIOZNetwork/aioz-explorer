package repository

import (
	"aioz.io/go-aioz/x_gob_explorer/context"
	"aioz.io/go-aioz/x_gob_explorer/domain/db"
	"aioz.io/go-aioz/x_gob_explorer/domain/entity"
	"fmt"
	"gorm.io/gorm"
	"log"
)

type TransactionRepo struct {
	db db.Database
}

func NewTransactionRepo(db db.Database) TransactionRepo {
	return TransactionRepo{
		db: db,
	}
}

func (tr *TransactionRepo) CreateTransactions(ctx context.Context, txns []*entity.Transaction) error {
	return tr.db.PrepareTransaction(
		func(txn *gorm.DB) error {
			for _, tx := range txns {
				if err := txn.Session(&gorm.Session{PrepareStmt: true}).Create(tx).Error; err != nil {
					return err
				}
			}
			return nil
		})
}

func (tr *TransactionRepo) MultiRowsInsert(ctx context.Context, txns []*entity.Transaction) error {
	if len(txns) > 0 {
		sql := `INSERT INTO transactions(block_height,block_time,hash,tx_index,gas,memo,fee,signatures,is_valid,payload,payload_log) VALUES`
		for i, t := range txns {
			sql += fmt.Sprintf(`(%v,%v,'%v',%v,%v,'%v','%v','%v',%v,'%v','%v')`,
				t.BlockHeight, t.BlockTime, t.Hash, t.TxIndex, t.Gas, t.Memo, t.Fee, t.Signatures, t.IsValid, t.Payload, t.PayloadLog)
			if i == len(txns)-1 {
				sql += ";"
			} else {
				sql += ","
			}
		}

		return tr.db.GetGormClient().Transaction(
			func(tx *gorm.DB) error {
				if err := tx.Exec(sql).Error; err != nil {
					log.Println(err)
				}
				return nil
			})
	}
	return nil
}
