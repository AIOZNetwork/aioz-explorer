package repository

import (
	"aioz.io/go-aioz/x_gob_explorer/context"
	"aioz.io/go-aioz/x_gob_explorer/domain/db"
	"aioz.io/go-aioz/x_gob_explorer/domain/entity"
	"fmt"
	"gorm.io/gorm"
	"log"
)

type MessageRepo struct {
	db db.Database
}

func NewMessageRepo(db db.Database) MessageRepo {
	return MessageRepo{
		db: db,
	}
}

func (mr *MessageRepo) CreateMessage(ctx context.Context, msg interface{}) error {
	return mr.db.PrepareTransaction(
		func(tx *gorm.DB) error {
			return tx.Session(&gorm.Session{PrepareStmt: true}).Create(msg).Error
		})
}

func (mr *MessageRepo) MultiRowsInsertMessages(ctx context.Context, msgs []*entity.Message) error {
	if len(msgs) > 0 {
		sql := `INSERT INTO messages(block_height,transaction_hash,message_index,message_time,type,is_valid,payload,payload_log) VALUES`
		for i, t := range msgs {
			sql += fmt.Sprintf(`(%v,'%v',%v,%v,'%v',%v,'%v','%v')`,
				t.BlockHeight, t.TransactionHash, t.MessageIndex, t.MessageTime, t.Type, t.IsValid, t.Payload, t.PayloadLog)
			if i == len(msgs)-1 {
				sql += ";"
			} else {
				sql += ","
			}
		}

		return mr.db.GetGormClient().Transaction(
			func(tx *gorm.DB) error {
				if err := tx.Exec(sql).Error; err != nil {
					log.Println(err)
				}
				return nil
			})
	}
	return nil
}

func (mr *MessageRepo) MultiRowsInsertMessageSend(ctx context.Context, msgs []*entity.MessageSend) error {
	if len(msgs) > 0 {
		sql := `INSERT INTO message_sends(block_height,transaction_hash,transaction_index,message_index,message_time,type,is_valid,payload,payload_log,from_address,to_address,amount) VALUES`
		for i, t := range msgs {
			sql += fmt.Sprintf(`(%v,'%v',%v,%v,%v,'%v',%v,'%v','%v','%v','%v','%v')`,
				t.BlockHeight, t.TransactionHash, t.TransactionIndex, t.MessageIndex, t.MessageTime, t.Type, t.IsValid, t.Payload, t.PayloadLog, t.FromAddress, t.ToAddress, t.Amount)
			if i == len(msgs)-1 {
				sql += ";"
			} else {
				sql += ","
			}
		}

		return mr.db.GetGormClient().Transaction(
			func(tx *gorm.DB) error {
				if err := tx.Exec(sql).Error; err != nil {
					log.Println(err)
				}
				return nil
			})
	}
	return nil
}
