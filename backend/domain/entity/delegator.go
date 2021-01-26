package entity

type Delegator struct {
	DelegatorAddress string `gorm:"primary_key;index:idx_delegator_delegator_address"`
	ValidatorAddress string `gorm:"primary_key;index:idx_delegator_validator_address"`
	WithdrawAddress  string
	Shares           string
}
