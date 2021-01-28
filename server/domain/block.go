package domain

import (
	"github.com/shopspring/decimal"
	"swagger-server/context"
)

type Block struct {
	ChainId               string
	Height                int64 `gorm:"primary_key"`
	AppHash               string
	HeaderHash            string
	ConsensusHash         string
	DataHash              string
	TxnsHash              string
	EvidenceHash          string
	LastCommitHash        string
	LastResultsHash       string
	NextValidatorsHash    string
	ProposerAddress       string
	ValidatorsHash        string
	NumTxs                int64
	Time                  int64
	TotalTxs              int64
	Size                  int
	AvgFee                string
	LastBlockId           string
	AnnualInflationRate   decimal.Decimal
	AnnualInflationBlocks decimal.Decimal
	InflationPerBlock     decimal.Decimal
	TotalSupply           string
	CirculatingSupply     string
}

type BlockRepository interface {
	GetByHeight(ctx context.Context, id string) (*Block, error)
	GetLatestBlocks(ctx context.Context, from, limit int) ([]*Block, error)
	GetByHash(ctx context.Context, hash string) (*Block, error)
	CountWithCondition(ctx context.Context, cond string, condParams ...interface{}) (int64, error)
	GetBestBlock(ctx context.Context) (*Block, error)
}

// Usecase represent the article's usecases
type BlockUsecase interface {
	GetByHeight(ctx context.Context, id string) (*Block, error)
	GetLatestBlocks(ctx context.Context, from, size int) ([]*Block, int64, error)
	GetByHash(ctx context.Context, hash string) (*Block, error)
	CountTotalBlocks(ctx context.Context) (int64, error)
	GetBestBlock(ctx context.Context) (*Block, error)
}
