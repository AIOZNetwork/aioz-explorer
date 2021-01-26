package controller

import (
	"aioz.io/go-aioz/x_gob_explorer/context"
	"aioz.io/go-aioz/x_gob_explorer/domain/entity"
	"aioz.io/go-aioz/x_gob_explorer/utils"
	"encoding/hex"
	"errors"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/tendermint/tendermint/rpc/core"
	"github.com/tendermint/tendermint/types"
)

func (ctrl *Controller) preprocessTransactionData(ctx context.Context,
	block *types.Block) ([]*entity.Transaction, string, error) {
	logger := ctx.Logger()

	txns := make([]*entity.Transaction, 0)
	listWallets := ""
	blockResults, err := core.BlockResults(nil, &block.Height)
	if err != nil {
		logger.Error(err)
		return nil, listWallets, err
	}
	for idx, tx := range block.Txs {
		logger.WithField("block", block.Hash().String())
		logger.WithField("txHash", tx.String())
		logger.WithField("txIndex", idx)
		txResult := blockResults.TxsResults[idx]
		sdkTx, err := txDecoder(tx)
		if err != nil {
			logger.Error(err)
			return nil, listWallets, err
		}
		stdTx, ok := sdkTx.(auth.StdTx)
		if !ok {
			logger.Error("Cannot parse transaction")
			return nil, listWallets, errors.New("cannot parse transaction")
		}
		fee, err1 := utils.JSONifyObject(stdTx.Fee)
		if err1 != nil {
			return nil, listWallets, err1
		}
		sigs, err1 := utils.JSONifyObject(stdTx.Signatures)
		jsonTx := cdc.MustMarshalJSON(stdTx)
		obj := &entity.Transaction{
			BlockHeight: block.Height,
			BlockTime:   block.Time.Unix(),
			Hash:        hex.EncodeToString(tx.Hash()),
			TxIndex:     int64(idx),
			Gas:         txResult.GasUsed,
			Memo:        stdTx.Memo,
			Fee:         string(fee),
			Signatures:  string(sigs),
			IsValid:     txResult.IsOK(),
			Payload:     string(jsonTx),
			PayloadLog:  txResult.Log,
		}
		txns = append(txns, obj)
		signers := stdTx.GetSigners()
		if len(signers) > 0 {
			listWallets = signers[0].String()
		}
	}
	return txns, listWallets, nil
}
