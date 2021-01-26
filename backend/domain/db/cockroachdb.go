package db

import (
	"aioz.io/go-aioz/x_gob_explorer/domain/entity"
	"aioz.io/go-aioz/x_gob_explorer/email"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"math"
	"math/rand"
	"os"
	"time"
)

var _ Database = &Cockroachdb{}

const maxRetries = 3

type txnFunc func(*gorm.DB) error

type Cockroachdb struct {
	dialect gorm.Dialector
	url     string
	client  *gorm.DB
	txns    []txnFunc
}

func NewCockroachDB(cruser, passwd, host, port, dbname, sslmode,
	sslrootcert, sslkey, sslcert string) Database {
	url := ""
	if sslmode == "disable" {
		url = fmt.Sprintf("user=%v host=%v port=%v dbname=%v",
			cruser, host, port, dbname)
	} else if sslmode == "require" {
		url = fmt.Sprintf("user=%v password=%v host=%v port=%v dbname=%v sslmode=%v sslrootcert=%v sslcert=%v sslkey=%v",
			cruser, passwd, host, port, dbname, sslmode, sslrootcert, sslcert, sslkey)
	} else {
		panic(errors.New("sslmode is undefined"))
	}
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: 200 * time.Second, // Slow SQL threshold
			LogLevel:      logger.Silent,
		},
	)
	c, err := gorm.Open(postgres.Open(url), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		PrepareStmt:                              true,
		SkipDefaultTransaction:                   true,
		Logger:                                   newLogger,
	})
	if err != nil {
		_ = email.SendMail(err.Error())
		panic(err)
	}

	db := &Cockroachdb{
		dialect: postgres.Open(url),
		url:     url,
		client:  c,
		txns:    make([]txnFunc, 0),
	}

	return db
}

func (db *Cockroachdb) GetGormClient() *gorm.DB {
	return db.client
}

func (db *Cockroachdb) PrepareTransaction(fn txnFunc) error {
	db.txns = append(db.txns, fn)
	return nil
}

func (db *Cockroachdb) FlushAllTxns() error {
	for _, fn := range db.txns {
		for retries := 0; retries < maxRetries; retries++ {
			if retries == maxRetries {
				return fmt.Errorf("hit max of %d retries, aborting", retries)
			}
			txn := db.client.Begin()
			if err := fn(txn); err != nil {
				// We need to cast GORM's db.Error to *pq.Error so we can
				// detect the Postgres transaction retry error code and
				// handle retries appropriately.
				switch errType := err.(type) {
				case *pq.Error:
					if errType.Code == "40001" {
						// Since this is a transaction retry error, we
						// ROLLBACK the transaction and sleep a little before
						// trying again.  Each time through the loop we sleep
						// for a little longer than the last time
						// (A.K.A. exponential backoff).
						txn.Rollback()
						var sleepMs = math.Pow(2, float64(retries)) * 100 * (rand.Float64() + 0.5)
						fmt.Printf("Hit 40001 transaction retry error, sleeping %s milliseconds\n", sleepMs)
						time.Sleep(time.Millisecond * time.Duration(10))
					} else {
						// If it's not a retry error, it's some other sort of
						// DB interaction error that needs to be handled by
						// the caller.
						return err
					}
				case *pgconn.PgError:
					if errType.Code == "40001" {
						// Since this is a transaction retry error, we
						// ROLLBACK the transaction and sleep a little before
						// trying again.  Each time through the loop we sleep
						// for a little longer than the last time
						// (A.K.A. exponential backoff).
						txn.Rollback()
						var sleepMs = math.Pow(2, float64(retries)) * 100 * (rand.Float64() + 0.5)
						fmt.Printf("Hit 40001 transaction retry error, sleeping %s milliseconds\n", sleepMs)
						time.Sleep(time.Millisecond * time.Duration(10))
					} else {
						// If it's not a retry error, it's some other sort of
						// DB interaction error that needs to be handled by
						// the caller.
						return err
					}
				default:
					return err
				}
			} else {
				// All went well, so we try to commit and break out of the
				// retry loop if possible.
				if err := txn.Commit().Error; err != nil {
					switch errType := err.(type) {
					case *pq.Error:
						if errType.Code == "40001" {
							// However, our attempt to COMMIT could also
							// result in a retry error, in which case we
							// continue back through the loop and try again.
							continue
						} else {
							// If it's not a retry error, it's some other sort
							// of DB interaction error that needs to be
							// handled by the caller.
							return err
						}
					case *pgconn.PgError:
						if errType.Code == "40001" {
							// However, our attempt to COMMIT could also
							// result in a retry error, in which case we
							// continue back through the loop and try again.
							continue
						} else {
							// If it's not a retry error, it's some other sort
							// of DB interaction error that needs to be
							// handled by the caller.
							return err
						}
					default:
						return err
					}
				}
				break
			}
		}
	}
	db.txns = make([]txnFunc, 0)
	return nil
}

func (db *Cockroachdb) Reset(tableSchema string) {
	// delete all tables
	if err := db.client.Transaction(func(tx *gorm.DB) error {
		return tx.Migrator().DropTable(
			&entity.Block{}, &entity.Transaction{}, &entity.Txs{}, &entity.MessageSend{},
			&entity.WalletAddress{}, &entity.Validator{}, &entity.Delegator{}, &entity.Stake{}, &entity.NodeInfo{}, &entity.PnTokenDevice{})
	}); err != nil {
		_ = email.SendMail(err.Error())
		panic(err)
	}

}
