package usecase

import (
	"swagger-server/context"
	"swagger-server/domain"
)

type deviceUsecase struct {
	deviceRepo domain.DeviceRepository
}

func NewDeviceUsecase(d domain.DeviceRepository) domain.DeviceUsecase {
	return &deviceUsecase{
		deviceRepo: d,
	}
}

func (du *deviceUsecase) SaveDevice(ctx context.Context, req *domain.SetPNTokenReq) error {
	return du.deviceRepo.SaveDevice(ctx, req)
}

func (du *deviceUsecase) DeleteDevice(ctx context.Context, deviceToken *domain.RemovePNTokenReq) error {
	return du.deviceRepo.DeleteDevice(ctx, deviceToken)
}

func (du *deviceUsecase) GetDevices(ctx context.Context) ([]*domain.PnTokenDevice, error) {
	return du.deviceRepo.GetDevices(ctx)
}

func (du *deviceUsecase) UpdateDevice(ctx context.Context, deviceToken string, status domain.DeviceStatus) error {
	return du.deviceRepo.UpdateDevice(ctx, deviceToken, status)
}

func (du *deviceUsecase) GetWalletsByToken(ctx context.Context, token string) ([]string, error) {
	wallets := make([]string, 0)
	devices, err := du.deviceRepo.GetWalletsByToken(ctx, token)
	for _, d := range devices {
		wallets = append(wallets, d.Wallet)
	}
	return wallets, err
}

func (du *deviceUsecase) GetTokensWallet(ctx context.Context, wallet string) ([]string, error) {
	tokens := make([]string, 0)
	devices, err := du.deviceRepo.GetTokensWallet(ctx, wallet)
	for _, d := range devices {
		tokens = append(tokens, d.PnToken)
	}
	return tokens, err
}
