package domain

import "swagger-server/context"

type StakingRepository interface {
	GetTopValidators(ctx context.Context, limit, offset int) ([]*Validator, error)
	GetTopStakingWallets(ctx context.Context, limit, offset int) ([]*Stake, error)
	GetTotalStakes(ctx context.Context) ([]*Stake, error)
	GetTotalStakedValidators(ctx context.Context) ([]*Validator, error)
	CountByTableNameWithCondition(ctx context.Context, table, cond string, condParams ...interface{}) (int64, error)
}

type StakingUsecase interface {
	GetTopValidators(ctx context.Context, limit, offset int) (*StakingValidatorStatistic, int64, error)
	GetTopStakingWallets(ctx context.Context, limit, offset int) (*StakingWalletStatisic, int64, error)
	GetTotalStakes(ctx context.Context) (string, error)
	GetTotalStakedValidators(ctx context.Context) (string, error)
}

type StakingValidatorStatistic struct {
	Validators  []*Validator
	TotalTokens string
}

type StakingWalletStatisic struct {
	Delegators  []*Stake
	TotalTokens string
}
