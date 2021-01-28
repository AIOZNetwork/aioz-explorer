package domain

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Validator struct {
	Address  string `gorm:"primary_key" json:"address"`
	Tokens   string `json:"tokens"`
	Power    int64  `json:"power"`
	Jailed   bool   `json:"jailed"`
	Status   string `json:"status"`
	IsActive bool   `json:"is_active"`

	Detail   string `json:"detail"`
	Identity string `json:"identity"`
	Moniker  string `json:"moniker"`
	Website  string `json:"website"`

	Period     uint64 `json:"period"`
	RewardPool string `json:"reward_pool"`

	ValConsAddr   string `json:"val_cons_addr"`
	ValConsPubkey string `json:"val_cons_pubkey"`
}

type ValidatorResp struct {
	Address    string       `gorm:"primary_key" json:"address"`
	Tokens     string       `json:"tokens"`
	Power      int64        `json:"power"`
	Jailed     bool         `json:"jailed"`
	Status     string       `json:"status"`
	StakedInfo []*Delegator `json:"staked_info"`
	IsActive   bool         `json:"is_active"`

	Detail   string `json:"detail"`
	Identity string `json:"identity"`
	Moniker  string `json:"moniker"`
	Website  string `json:"website"`

	Period     uint64    `json:"period"`
	RewardPool sdk.DecCoins `json:"reward_pool"`

	ValConsAddr   string `json:"val_cons_addr"`
	ValConsPubkey string `json:"val_cons_pubkey"`
}

type ValidatorRepository interface {
	GetValidatorByValAddress(ctx context.Context, address string) (*Validator, error)
}

type ValidatorUsecase interface {
	GetValidatorByValAddress(ctx context.Context, address string) (*Validator, error)
}
