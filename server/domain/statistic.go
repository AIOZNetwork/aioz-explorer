package domain

import "swagger-server/context"

type Statistic struct {
	TotalBlock        int64
	TotalTransaction  int64
	TotalWallets      int64
	Volume24h         string
	AvgFee            string
	TopAccount        []*WalletResp
	Inflation         string
	MarketCap         string
	BlockTime         string
	CirculatingSupply string
	TotalSupply       string
}

type StatisticRepository interface {
}

type StatisticUsecase interface {
	GetStatistic(ctx context.Context) (*Statistic, error)
	GetInflation(ctx context.Context) (string, error)
	GetTotalSupply(ctx context.Context) (string, error)
}
