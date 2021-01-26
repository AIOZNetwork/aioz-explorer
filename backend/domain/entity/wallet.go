package entity

import (
	"github.com/shopspring/decimal"
)

type WalletAddress struct {
	Address       string `gorm:"primary_key"`
	Coins         string
	AccountNumber uint64
	Sequence      uint64
	PubKey        string
	CoinAIOZ      decimal.Decimal `gorm:"type:numeric;index:idx_wallet_address_coin_aioz"`
	CoinStake     decimal.Decimal `gorm:"type:numeric;index:idx_wallet_address_coin_stake"`

	OriginalVesting  string
	DelegatedFree    string
	DelegatedVesting string
	StartTime        int64 `gorm:"index:idx_wallet_address_start_time"`
	EndTime          int64 `gorm:"index:idx_wallet_address_end_time"`
}
