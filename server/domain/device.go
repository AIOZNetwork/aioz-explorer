package domain

import "swagger-server/context"

type PushTokenRequest struct {
	Token   string
	Wallets []string
}

type SetPNTokenReq struct {
	PlatformName Platform `json:"platform"`
	PnToken      string   `json:"token"`
	Wallets      []string `json:"wallets"`
}

type RemovePNTokenReq struct {
	PnToken string   `json:"token"`
	Wallets []string `json:"wallets"`
}

type PnTokenDevice struct {
	PlatformName string
	PnToken      string `gorm:"PRIMARY_KEY;NOT NULL;index:idx_device_token"`
	Wallet       string `gorm:"PRIMARY_KEY;NOT NULL;index:idx_device_wallet"`
	CreatedAt    int64
	DeviceStatus DeviceStatus
}

type DeviceUsecase interface {
	SaveDevice(ctx context.Context, req *SetPNTokenReq) error
	DeleteDevice(ctx context.Context, deviceToken *RemovePNTokenReq) error
	GetDevices(ctx context.Context) ([]*PnTokenDevice, error)
	UpdateDevice(ctx context.Context, deviceToken string, status DeviceStatus) error
	GetWalletsByToken(ctx context.Context, token string) ([]string, error)
	GetTokensWallet(ctx context.Context, wallet string) ([]string, error)
}

type DeviceRepository interface {
	SaveDevice(ctx context.Context, req *SetPNTokenReq) error
	DeleteDevice(ctx context.Context, deviceToken *RemovePNTokenReq) error
	GetDevices(ctx context.Context) ([]*PnTokenDevice, error)
	UpdateDevice(ctx context.Context, deviceToken string, status DeviceStatus) error
	GetWalletsByToken(ctx context.Context, token string) ([]*PnTokenDevice, error)
	GetTokensWallet(ctx context.Context, wallet string) ([]*PnTokenDevice, error)
}
