package entity

import "github.com/shopspring/decimal"

type Block struct {
	ChainId               string
	Height                int64 `gorm:"primary_key;index:idx_block_height"`
	AppHash               string
	HeaderHash            string `gorm:"index:idx_block_header_hash"`
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
	AnnualInflationRate   decimal.Decimal `gorm:"type:numeric;index:idx_block_annual_rate"`
	AnnualInflationBlocks decimal.Decimal `gorm:"type:numeric;index:idx_block_annual_blocks"`
	InflationPerBlock     decimal.Decimal `gorm:"type:numeric;index:idx_block_per_block"`
	TotalSupply           string
	CirculatingSupply     string
}
