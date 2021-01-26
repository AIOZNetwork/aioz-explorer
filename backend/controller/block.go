package controller

import (
	"aioz.io/go-aioz/x_gob_explorer/context"
	"aioz.io/go-aioz/x_gob_explorer/domain/entity"
	"aioz.io/go-aioz/x_gob_explorer/utils"
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/shopspring/decimal"
	"github.com/tendermint/tendermint/rpc/core"
	"github.com/tendermint/tendermint/types"
	"log"
)

var cummTxs = int64(0)

func (ctrl *Controller) preprocessBlockData(ctx context.Context, block *types.Block) *entity.Block {

	prevBlock, err := ctrl.blockRepo.GetPreviousBlockInfo()
	if err != nil {
		cummTxs = 0
	} else {
		cummTxs = prevBlock.TotalTxs
	}
	numTxs := int64(len(block.Txs))
	cummTxs += numTxs

	totalAmount := sdk.NewCoins()
	for _, tx := range block.Txs {
		sdkTx, err := txDecoder(tx)
		if err != nil {
			return nil
		}
		stdTx, ok := sdkTx.(auth.StdTx)
		if !ok {
			return nil
		}
		totalAmount = totalAmount.Add(stdTx.Fee.Amount...)
	}
	for i, c := range totalAmount {
		amt := c.Amount.Quo(sdk.NewInt(numTxs))
		totalAmount[i] = sdk.NewCoin(c.Denom, amt)
	}
	blockResults, err := core.BlockResults(nil, &block.Height)
	if err != nil {
		return nil
	}
	annualInflationBlocks := ""
	annualInflationRate := ""
	inflationPerBlock := ""
	for _, e := range blockResults.BeginBlockEvents {
		switch e.Type {
		case "mint":
			for _, attr := range e.Attributes {
				switch string(attr.Key) {
				case "inflation":
					annualInflationRate = string(attr.Value)
				case "annual_provisions":
					annualInflationBlocks = string(attr.Value)
				case "amount":
					inflationPerBlock = string(attr.Value)
				}
			}
		}
	}
	aib, _ := decimal.NewFromString(annualInflationBlocks)
	air, _ := decimal.NewFromString(annualInflationRate)
	ipb, _ := decimal.NewFromString(inflationPerBlock)

	// calculate spendable coins after this block
	wallets, err := ctrl.walletRepo.GetVestingWallet(ctx, block.Time)
	if err != nil {
		log.Println(err)
	}
	lockedCoins := sdk.NewCoins()
	totalSupply := supplyKeeper.GetSupply(ctx.CMCtx()).GetTotal()
	//totalSupply, err := ctrl.lcdRepo.GetTotalSupply()
	//if err != nil {
	//	log.Println(err)
	//	totalSupply, _ = sdk.ParseCoins(prevBlock.TotalSupply)
	//}
	for _, w := range wallets {
		accAddr, _ := sdk.AccAddressFromBech32(w.Address)
		acc := accountKeeper.GetAccount(ctx.CMCtx(), accAddr)
		lockedCoins = lockedCoins.Add(acc.SpendableCoins(block.Time)...)
	}
	cirSupply := totalSupply.Sub(lockedCoins)
	b, _ := json.Marshal(totalAmount)
	obj := &entity.Block{
		ChainId:            block.ChainID,
		Height:             block.Height,
		AppHash:            block.AppHash.String(),
		HeaderHash:         block.Header.Hash().String(),
		ConsensusHash:      block.ConsensusHash.String(),
		DataHash:           block.DataHash.String(),
		TxnsHash:           block.Data.Hash().String(),
		LastCommitHash:     block.LastCommitHash.String(),
		LastResultsHash:    block.LastResultsHash.String(),
		NextValidatorsHash: block.NextValidatorsHash.String(),
		ProposerAddress:    sdk.ConsAddress(block.ProposerAddress).String(),
		ValidatorsHash:     block.ValidatorsHash.String(),
		NumTxs:             int64(len(block.Txs)),
		Time:               utils.ConvertNanoTimestampToMilliSecond(block.Time.UnixNano()),
		TotalTxs:           cummTxs,
		EvidenceHash:       block.EvidenceHash.String(),
		Size:               block.Size(),
		AvgFee:             string(b),
		LastBlockId:        block.LastBlockID.Hash.String(),

		AnnualInflationBlocks: aib,
		AnnualInflationRate:   air,
		InflationPerBlock:     ipb,
		TotalSupply:           totalSupply.String(),
		CirculatingSupply:     cirSupply.String(),
	}

	return obj
}
