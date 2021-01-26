package entity

type Transaction struct {
	BlockHeight int64 `gorm:"primary_key;index:idx_transaction_block_height;index:idx_transaction_height_index,priority:1"`
	BlockTime   int64
	Hash        string `gorm:"index:idx_transaction_hash"`
	TxIndex     int64  `gorm:"primary_key;index:idx_transaction_tx_index;index:idx_transaction_height_index,priority:2"`
	Gas         int64
	Memo        string
	Fee         string
	Signatures  string
	IsValid     bool
	Payload     string
	PayloadLog  string
}
