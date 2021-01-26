package entity

type Validator struct {
	Address  string `gorm:"primary_key"`
	Tokens   string
	Power    int64
	Jailed   bool
	Status   string
	IsActive bool

	Detail   string
	Identity string
	Moniker  string
	Website  string

	Period     uint64
	RewardPool string

	ValConsAddr   string
	ValConsPubkey string
}
