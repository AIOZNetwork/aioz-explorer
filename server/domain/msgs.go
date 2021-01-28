package domain

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lib/pq"
	"swagger-server/context"
)

type Txs struct {
	//gorm.Model
	BlockHeight      int64          `json:"block_height"`
	BlockHash        string         `json:"block_hash"`
	BlockTime        int64          `json:"block_time"`
	TransactionHash  string         `json:"transaction_hash"`
	TransactionIndex int64          `json:"transaction_index"`
	MessageType      string         `json:"message_type"`
	MessageIndex     int64          `json:"message_index"`
	Address          string         `json:"address"`
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

type MsgsRepository interface {
	GetRecentMsgs(ctx context.Context, limit, offset int) ([]*Txs, error)
	GetMsgsByAddress(ctx context.Context, address string, limit, offset int) ([]*Txs, error)
	GetMsgSendHistory(ctx context.Context, heightFrom, heightTo, limit, offset int) ([]*MessageSend, error)
	GetVolumeSend24h(ctx context.Context) ([]*Txs, error)
	CountByTableNameWithCondition(ctx context.Context, table, cond string, condParams ...interface{}) (int64, error)

	GetMsgsByTxHash(ctx context.Context, txHash string) ([]*Txs, error)
	PaginateVolumeSend24h(ctx context.Context) ([]*Txs, error)
}

type MsgsUsecase interface {
	GetRecentMsgs(ctx context.Context, limit, offset int) ([]*Txs, int64, error)
	GetMsgsByAddress(ctx context.Context, address string, limit, offset int) ([]*Txs, int64, error)
	GetMsgSendHistory(ctx context.Context, heightFrom, heightTo, limit, offset int) ([]*MessageSend, int64, error)
	GetVolumeSend24h(ctx context.Context) (string, error)

	GetMsgsByTxHash(ctx context.Context, txHash string) ([]*Txs, error)
}
