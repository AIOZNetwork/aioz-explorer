package entity

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lib/pq"
)

type Txs struct {
	BlockHeight      int64          `gorm:"primary_key;index:idx_txs_block_height;index:idx_txs_comp_address,priority:2" json:"block_height"`
	BlockHash        string         `json:"block_hash"`
	BlockTime        int64          `json:"block_time" gorm:"index:idx_txs_block_time;index:idx_txs_send"`
	TransactionHash  string         `gorm:"index:idx_txs_transaction_hash" json:"transaction_hash"`
	TransactionIndex int64          `gorm:"primary_key;index:idx_txs_transaction_index;index:idx_txs_comp_address,priority:3" json:"transaction_index"`
	MessageType      string         `json:"message_type" gorm:"index:idx_txs_message_type;index:idx_txs_send"`
	MessageIndex     int64          `gorm:"primary_key;index:idx_txs_message_index;index:idx_txs_comp_address,priority:4" json:"message_index"`
	Address          string         `gorm:"primary_key;index:idx_txs_address;index:idx_txs_comp_address,priority:1" json:"address"`
	AddressFrom      pq.StringArray `gorm:"type:text[]" json:"address_from"`
	AddressTo        pq.StringArray `gorm:"type:text[]" json:"address_to"`
	Amount           string         `json:"amount"`
	IsValid          bool           `json:"is_valid"`
	PayloadMsg       string         `json:"payload_msg"`
	PayloadErr       string         `json:"payload_err"`
}

type TxsResp struct {
	BlockHeight      int64          `json:"block_height"`
	BlockHash        string         `json:"block_hash"`
	BlockTime        int64          `json:"block_time"`
	TransactionHash  string         `json:"transaction_hash"`
	TransactionIndex int64          `json:"transaction_index"`
	TransactionType  string         `json:"transaction_type"`
	MessageType      string         `json:"message_type"`
	MessageIndex     int64          `json:"message_index"`
	Address          string         `json:"address"`
	AddressFrom      pq.StringArray `json:"address_from"`
	AddressTo        pq.StringArray `json:"address_to"`
	Amount           sdk.Coins      `json:"amount"`
	IsValid          bool           `json:"is_valid"`
	PayloadMsg       string         `json:"payload_msg"`
	PayloadErr       string         `json:"payload_err"`
}
