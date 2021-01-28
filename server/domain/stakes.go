package domain

import "github.com/shopspring/decimal"

type Stake struct {
	DelegatorAddress string `gorm:"primary_key" json:"delegator_address"`
	Shares           string `json:"shares"`
	SharesDec        decimal.Decimal `gorm:"type:numeric;index:idx_staked_share_dec"`
}
