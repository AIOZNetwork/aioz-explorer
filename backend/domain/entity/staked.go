package entity

import "github.com/shopspring/decimal"

type Stake struct {
	DelegatorAddress string `gorm:"primary_key"`
	Shares           string
	SharesDec        decimal.Decimal `gorm:"type:numeric;index:idx_staked_share_dec"`
}
