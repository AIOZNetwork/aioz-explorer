package domain

import (
	"database/sql"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"swagger-server/context"
)

type WalletAddress struct {
	Address       string `gorm:"primary_key"`
	Coins         string
	AccountNumber uint64
	Sequence      uint64
	PubKey        string

	OriginalVesting  string
	DelegatedFree    string
	DelegatedVesting string
	StartTime        int64
	EndTime          int64
}

type WalletResp struct {
	Address       string    `json:"address"`
	Coins         sdk.Coins `json:"coins"`
	AccountNumber uint64    `json:"account_number"`
	Sequence      uint64    `json:"sequence"`
	PubKey        string    `json:"pub_key"`
	StakedInfo    []*Delegator

	OriginalVesting  string `json:"original_vesting"`
	DelegatedFree    string `json:"delegated_free"`
	DelegatedVesting string `json:"delegated_vesting"`
	StartTime        int64  `json:"start_time"`
	EndTime          int64  `json:"end_time"`
}

type RequestNewKey struct {
	Password string `json:"password,omitempty"`
	Account  uint32 `json:"account"`
	Index    uint32 `json:"index"`
}

type RequestRecoverKey struct {
	Mnemonic string `json:"mnemonic"`
	Password string `json:"password"`
	Account  uint32 `json:"account"`
	Index    uint32 `json:"index"`
}

type KeyResponse struct {
	Address   string `json:"address,omitempty"`
	PubKey    []byte `json:"pub_key,omitempty"`
	PrivKey   string `json:"priv_key,omitempty"`
	Mnemonic  string `json:"mnemonic,omitempty"`
	PrivArmor string `json:"priv_armor,omitempty"`
}

type RequestEncrKey struct {
	PrivKey  string `json:"priv_key"`
	Password string `json:"password"`
}

type RequestDecrKey struct {
	PrivArmor string `json:"priv_armor"`
	Password  string `json:"password"`
}

type TxSend struct {
	PrivKey   string `json:"priv_key"`
	ToAddress string `json:"to_address"`
	Amount    string `json:"amount"`
}

type ContactsWallet struct {
	Address        string    `json:"address"`
	TotalTxs       int64     `json:"total_txs"`
	AmountSent     sdk.Coins `json:"amount_sent"`
	AmountReceived sdk.Coins `json:"amount_received"`
}

type WalletRepository interface {
	GetWalletByAddress(ctx context.Context, address string) (*WalletAddress, error)
	GetTxsByAddress(ctx context.Context, address string, limit, offset int) ([]*MessageSend, error)
	GetContactsByWallet(ctx context.Context, address string, condition string) (*sql.Rows, error)
	GetTopWallet(ctx context.Context, topN int64) ([]*WalletAddress, error)
	CountByTableNameWithCondition(ctx context.Context, table, cond string, condParams ...interface{}) (int64, error)
	GetVolumeFromWalletToWallet(ctx context.Context, from, to string) ([]*MessageSend, error)
}

type WalletUsecase interface {
	GetWalletByAddress(ctx context.Context, address string) (*WalletResp, error)
	GetValidatorByAddress(ctx context.Context, address string) (*ValidatorResp, error)
	GetTxsByAddress(ctx context.Context, address string, limit, offset int) ([]*MessageSend, int64, error)
	GetTxsByAddressV2(ctx context.Context, address string, limit, offset int) ([]*TxsResp, int64, error)
	CreateWallet(ctx context.Context, req *RequestNewKey) (*KeyResponse, error)
	RecoverWallet(ctx context.Context, req *RequestRecoverKey) (*KeyResponse, error)
	GetContactsWallet(ctx context.Context, addresses []string) ([]*ContactsWallet, error)
	EncryptKey(ctx context.Context, req *RequestEncrKey) (*KeyResponse, error)
	DecryptKey(ctx context.Context, req *RequestDecrKey) (*KeyResponse, error)
	GetTopWallet(ctx context.Context, topN int64) ([]*WalletResp, error)
	CountTotalWallets(ctx context.Context) (int64, error)
	GetVolumeFromWalletToWallet(ctx context.Context, from, to string) (sdk.Coins, error)
}
