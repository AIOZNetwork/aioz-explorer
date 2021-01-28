package usecase

import (
	"github.com/cosmos/cosmos-sdk/types"
	"log"
	"swagger-server/context"
	"swagger-server/domain"
)

type stakingUsecase struct {
	stakingRepo domain.StakingRepository
}

func NewStakingUsecase(t domain.StakingRepository) domain.StakingUsecase {
	return &stakingUsecase{
		stakingRepo: t,
	}
}

func (t *stakingUsecase) GetTopValidators(ctx context.Context, limit, offset int) (*domain.StakingValidatorStatistic, int64, error) {
	resp, err := t.stakingRepo.GetTopValidators(ctx, limit, offset)
	if err != nil {
		return nil, -1, err
	}
	total, err := t.stakingRepo.CountByTableNameWithCondition(ctx, domain.Table_validator, "")
	if err != nil {
		return nil, -1, err
	}
	totalTokens, err := t.GetTotalStakedValidators(ctx)
	if err != nil {
		return nil, -1, err
	}
	//totalTokens := decimal.NewFromFloat(0)
	//for _, t := range resp {
	//	d, err := decimal.NewFromString(t.Tokens)
	//	if err != nil {
	//		return nil, -1, err
	//	}
	//	totalTokens = totalTokens.Add(d)
	//}

	ret := &domain.StakingValidatorStatistic{
		Validators:  resp,
		TotalTokens: totalTokens,
	}

	return ret, total, nil
}

func (t *stakingUsecase) GetTopStakingWallets(ctx context.Context, limit, offset int) (*domain.StakingWalletStatisic, int64, error) {
	resp, err := t.stakingRepo.GetTopStakingWallets(ctx, limit, offset)
	if err != nil {
		return nil, -1, err
	}
	total, err := t.stakingRepo.CountByTableNameWithCondition(ctx, domain.Table_stake, "")
	if err != nil {
		return nil, -1, err
	}
	totalShares, err := t.GetTotalStakes(ctx)
	if err != nil {
		return nil, -1, err
	}
	//totalShares := decimal.NewFromFloat(0)
	//for _, r := range resp {
	//	totalShares = totalShares.Add(r.SharesDec)
	//}

	ret := &domain.StakingWalletStatisic{
		Delegators:  resp,
		TotalTokens: totalShares,
	}
	return ret, total, nil
}

func (t *stakingUsecase) GetTotalStakes(ctx context.Context) (string, error) {
	result, err := t.stakingRepo.GetTotalStakes(ctx)
	if err != nil {
		return "0", err
	}
	if len(result) == 0 {
		return "0", nil
	}
	total, err := types.NewDecFromStr("0")
	if err != nil {
		return "0", err
	}
	for _, v := range result {
		temp, err := types.NewDecFromStr(v.Shares)
		if err != nil {
			log.Println(err)
			continue
		}
		total = total.Add(temp)
	}
	return total.String(), nil
}

func (t *stakingUsecase) GetTotalStakedValidators(ctx context.Context) (string, error) {
	result, err := t.stakingRepo.GetTotalStakedValidators(ctx)
	if err != nil {
		return "0", err
	}
	if len(result) == 0 {
		return "0", nil
	}
	total, err := types.NewDecFromStr("0")
	if err != nil {
		return "0", err
	}
	for _, v := range result {
		temp, err := types.NewDecFromStr(v.Tokens)
		if err != nil {
			log.Println(err)
			continue
		}
		total = total.Add(temp)
	}
	return total.String(), nil
}
