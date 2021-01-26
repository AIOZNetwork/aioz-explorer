package controller

import (
	"aioz.io/go-aioz/x_gob_explorer/context"
	"aioz.io/go-aioz/x_gob_explorer/domain/entity"
	"errors"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/shopspring/decimal"
)

func (ctrl *Controller) updateStakingDelegate(ctx context.Context, msg staking.MsgDelegate) error {
	validator, found := stakingKeeper.GetValidator(ctx.CMCtx(), msg.ValidatorAddress)
	if !found {
		return errors.New("no validator found")
	}

	objVal := &entity.Validator{
		Address: validator.OperatorAddress.String(),
		Tokens:  validator.Tokens.String(),
		Power:   validator.GetConsensusPower(),
		Jailed:  validator.Jailed,
		Status:  validator.Status.String(),
	}
	if err := ctrl.validatorRepo.UpsertValidator(ctx, objVal); err != nil {
		return err
	}

	delegator, found := stakingKeeper.GetDelegation(ctx.CMCtx(), msg.DelegatorAddress, msg.ValidatorAddress)
	if !found {
		return errors.New("no validator found")
	}
	objDel := &entity.Delegator{
		DelegatorAddress: delegator.DelegatorAddress.String(),
		ValidatorAddress: delegator.ValidatorAddress.String(),
		Shares:           delegator.Shares.String(),
	}
	if err := ctrl.delegatorRepo.UpsertDelegator(ctx, objDel); err != nil {
		return err
	}

	stakeAcc, err := ctrl.stakeRepo.GetCurrentStakeByDelegator(ctx, msg.DelegatorAddress.String())
	if err != nil {
		return err
	}
	dec, err := decimal.NewFromString(msg.Amount.Amount.String())
	if err != nil {
		return err
	}
	if stakeAcc == nil {
		objStake := &entity.Stake{
			DelegatorAddress: msg.DelegatorAddress.String(),
			Shares:           msg.Amount.Amount.String(),
			SharesDec:        dec,
		}
		return ctrl.stakeRepo.UpsertStakedCoins(ctx, objStake)
	}
	stakeAcc.SharesDec = stakeAcc.SharesDec.Add(dec)
	stakeAcc.Shares = stakeAcc.SharesDec.String()
	return ctrl.stakeRepo.UpsertStakedCoins(ctx, stakeAcc)
}

func (ctrl *Controller) updateStakingUndelegate(ctx context.Context, msg staking.MsgUndelegate) error {
	validator, found := stakingKeeper.GetValidator(ctx.CMCtx(), msg.ValidatorAddress)
	if !found {
		return errors.New("no validator found")
	}

	objVal := &entity.Validator{
		Address: validator.OperatorAddress.String(),
		Tokens:  validator.Tokens.String(),
		Power:   validator.GetConsensusPower(),
		Jailed:  validator.Jailed,
		Status:  validator.Status.String(),
	}
	if err := ctrl.validatorRepo.UpsertValidator(ctx, objVal); err != nil {
		return err
	}

	delegator, found := stakingKeeper.GetDelegation(ctx.CMCtx(), msg.DelegatorAddress, msg.ValidatorAddress)
	if !found {
		return errors.New("no validator found")
	}
	objDel := &entity.Delegator{
		DelegatorAddress: delegator.DelegatorAddress.String(),
		ValidatorAddress: delegator.ValidatorAddress.String(),
		Shares:           delegator.Shares.String(),
	}
	return ctrl.delegatorRepo.UpsertDelegator(ctx, objDel)
}

func (ctrl *Controller) updateBeginRedelegate(ctx context.Context, msg staking.MsgBeginRedelegate) error {
	validator1, found := stakingKeeper.GetValidator(ctx.CMCtx(), msg.ValidatorSrcAddress)
	if !found {
		return errors.New("no validator found")
	}
	objVal1 := &entity.Validator{
		Address: validator1.OperatorAddress.String(),
		Tokens:  validator1.Tokens.String(),
		Power:   validator1.GetConsensusPower(),
		Jailed:  validator1.Jailed,
		Status:  validator1.Status.String(),
	}
	if err := ctrl.validatorRepo.UpsertValidator(ctx, objVal1); err != nil {
		return err
	}

	validator2, found := stakingKeeper.GetValidator(ctx.CMCtx(), msg.ValidatorDstAddress)
	if !found {
		return errors.New("no validator found")
	}
	objVal2 := &entity.Validator{
		Address: validator2.OperatorAddress.String(),
		Tokens:  validator2.Tokens.String(),
		Power:   validator2.GetConsensusPower(),
		Jailed:  validator2.Jailed,
		Status:  validator2.Status.String(),
	}
	if err := ctrl.validatorRepo.UpsertValidator(ctx, objVal2); err != nil {
		return err
	}

	delegator1, found := stakingKeeper.GetDelegation(ctx.CMCtx(), msg.DelegatorAddress, msg.ValidatorSrcAddress)
	if !found {
		return errors.New("no validator found")
	}
	objDel1 := &entity.Delegator{
		DelegatorAddress: delegator1.DelegatorAddress.String(),
		ValidatorAddress: delegator1.ValidatorAddress.String(),
		Shares:           delegator1.Shares.String(),
	}
	if err := ctrl.delegatorRepo.UpsertDelegator(ctx, objDel1); err != nil {
		return err
	}

	delegator2, found := stakingKeeper.GetDelegation(ctx.CMCtx(), msg.DelegatorAddress, msg.ValidatorDstAddress)
	if !found {
		return errors.New("no validator found")
	}
	objDel2 := &entity.Delegator{
		DelegatorAddress: delegator2.DelegatorAddress.String(),
		ValidatorAddress: delegator2.ValidatorAddress.String(),
		Shares:           delegator2.Shares.String(),
	}
	if err := ctrl.delegatorRepo.UpsertDelegator(ctx, objDel2); err != nil {
		return err
	}

	return nil
}

func (ctrl *Controller) updateMsgUnjail(ctx context.Context, msg slashing.MsgUnjail) error {
	validator, found := stakingKeeper.GetValidator(ctx.CMCtx(), msg.ValidatorAddr)
	if !found {
		return nil
	}
	obj := &entity.Validator{
		Address: validator.OperatorAddress.String(),
		Tokens:  validator.Tokens.String(),
		Power:   validator.GetConsensusPower(),
		Jailed:  validator.Jailed,
		Status:  validator.Status.String(),
	}
	return ctrl.validatorRepo.UpsertValidator(ctx, obj)
}
