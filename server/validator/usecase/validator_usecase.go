package usecase

import (
	"context"
	"swagger-server/domain"
)

type validatorUsecase struct {
	validatorRepo domain.ValidatorRepository
}

func NewValidatorUsecase(vr domain.ValidatorRepository) domain.ValidatorUsecase {
	return &validatorUsecase{
		validatorRepo: vr,
	}
}
/*
func (v *validatorUsecase) UpdateValidatorsFromNode(ctx context.Context, body *domain.Validator) error {
	return v.validatorRepo.UpdateValidatorsFromNode(ctx, body)
}

func (v *validatorUsecase) GetValidatorsInfo(ctx context.Context, limit, offset int64) ([]*domain.Validator, error) {
	return v.validatorRepo.GetValidatorsInfo(ctx, limit, offset)
}
*/
func (v *validatorUsecase) GetValidatorByValAddress(ctx context.Context, address string) (*domain.Validator, error) {
	return v.validatorRepo.GetValidatorByValAddress(ctx, address)
}
