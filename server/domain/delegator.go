package domain

import "context"

type DelegatorRepository interface {
	GetDelegatorByAccAddress(ctx context.Context, accAddress string) ([]*Delegator, error)
	GetDelegatorByValAddress(ctx context.Context, valAddress string) ([]*Delegator, error)
	CountWithCondition(ctx context.Context, cond string, condParams ...interface{}) (int64, error)
}

type DelegatorUsecase interface {
	GetDelegatorByAccAddress(ctx context.Context, accAddress string) ([]*Delegator, int64, error)
	GetDelegatorByValAddress(ctx context.Context, valAddress string) ([]*Delegator, int64, error)
}

type Delegator struct {
	DelegatorAddress string `gorm:"primary_key" json:"delegator_address"`
	ValidatorAddress string `gorm:"primary_key" json:"validator_address"`
	Shares           string `json:"shares"`
}
