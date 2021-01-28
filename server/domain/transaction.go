package domain

import "swagger-server/context"

type Transaction struct {
	BlockHeight int64
	BlockTime   int64
	Hash        string `gorm:"primary_key"`
	TxIndex     int64  `gorm:"index:idx_transaction_tx_index"`
	Gas         int64
	Memo        string
	Fee         string
	Signatures  string
	IsValid     bool
	Payload     string
}

type TransactionResp struct {
	BlockHeight int64
	BlockTime   int64
	Hash        string `gorm:"primary_key"`
	TxIndex     int64  `gorm:"index:idx_transaction_tx_index"`
	Gas         int64
	Memo        string
	Fee         string
	Signatures  string
	IsValid     bool
	Payload     []*Txs
}

type TransactionRepository interface {
	GetTransaction(ctx context.Context, hash string) (*Transaction, error)
	GetTransactions(ctx context.Context, limit, offset int64) ([]*Transaction, error)
	GetTransactionsByBlockHeight(ctx context.Context, height, limit, offset int64) ([]*Transaction, error)
	CountWithCondition(ctx context.Context, cond string, condParams ...interface{}) (int64, error)
}

type TransactionUsecase interface {
	GetTransaction(Ctx context.Context, hash string) (*TransactionResp, error)
	GetTransactions(Ctx context.Context, limit, offset int64) ([]*TransactionResp, int64, error)
	GetTransactionsByBlockHeight(ctx context.Context, height, limit, offset int64) ([]*TransactionResp, int64, error)
	CountTotalTx(ctx context.Context) (int64, error)
}
