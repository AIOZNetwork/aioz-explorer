package entity

type PnTokenDevice struct {
	PlatformName string
	PnToken      string `gorm:"PRIMARY_KEY;NOT NULL;index:idx_device_token"`
	Wallet       string `gorm:"PRIMARY_KEY;NOT NULL;index:idx_device_wallet"`
	CreatedAt    int64
	DeviceStatus DeviceStatus
}

type DeviceStatus int

const (
	Device_status_deactive = 0
	Device_status_active   = 1

	MessageNotiReceive = "recipient"
	MessageNotiSend    = "sender"

	DefaultDecCoin = "1000000000000000000"
)

