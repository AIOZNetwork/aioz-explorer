package controller

import (
	"aioz.io/go-aioz/x_gob_explorer/context"
	"aioz.io/go-aioz/x_gob_explorer/domain/entity"
)

func (ctrl *Controller) handleValidatorsPool(ctx context.Context) error {
	vals := stakingKeeper.GetAllValidators(ctx.CMCtx())
	for _, v := range vals {
		rewards := distributionKeeper.GetValidatorCurrentRewards(ctx.CMCtx(), v.OperatorAddress)
		objVal := &entity.Validator{
			Address:    v.OperatorAddress.String(),
			Period:     rewards.Period,
			RewardPool: rewards.Rewards.String(),
		}
		if err := ctrl.validatorRepo.UpsertValidatorRewardPool(ctx, objVal); err != nil {
			return err
		}
	}
	return nil
}
