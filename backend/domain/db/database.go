package db

import (
	// Import GORM-related packages.
	"gorm.io/gorm"
)

type Database interface {
	GetGormClient() *gorm.DB
	PrepareTransaction(fn txnFunc) error
	FlushAllTxns() error
	Reset(string)

}
