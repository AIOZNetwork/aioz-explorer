package usecase

import (
	"context"
	"swagger-server/domain"
)

type delegatorUsecase struct {
	delegatorRepo domain.DelegatorRepository
}

func NewDelegatorUsecase(d domain.DelegatorRepository) domain.DelegatorUsecase {
	return &delegatorUsecase{
		delegatorRepo: d,
	}
}

func (d *delegatorUsecase) GetDelegatorByAccAddress(ctx context.Context, address string) ([]*domain.Delegator, int64, error) {
	resp, err := d.delegatorRepo.GetDelegatorByAccAddress(ctx, address)
	if err != nil {
		return nil, -1, err
	}
	total, err := d.delegatorRepo.CountWithCondition(ctx, "delegator_address = ?", address)
	if err != nil {
		return nil, -1, err
	}
	return resp, total, nil
}

func (d *delegatorUsecase) GetDelegatorByValAddress(ctx context.Context, valAddress string) ([]*domain.Delegator, int64, error) {
	resp, err := d.delegatorRepo.GetDelegatorByValAddress(ctx, valAddress)
	if err != nil {
		return nil, -1, err
	}
	total, err := d.delegatorRepo.CountWithCondition(ctx, "validator_address = ?", valAddress)
	if err != nil {
		return nil, -1, err
	}
	return resp, total, nil
}