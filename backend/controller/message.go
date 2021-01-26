package controller

import (
	"aioz.io/go-aioz/x_gob_explorer/context"
	"aioz.io/go-aioz/x_gob_explorer/domain/entity"
	"aioz.io/go-aioz/x_gob_explorer/utils"
	"encoding/hex"
	"errors"
	"fmt"
	cmtypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/shopspring/decimal"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/rpc/core"
	tmtypes "github.com/tendermint/tendermint/types"
)

func (ctrl *Controller) preprocessMessageData(ctx context.Context,
	block *tmtypes.Block) (map[string]bool, []*entity.Txs, error) {
	mapMsgWithdrawRwdEvt := make(map[string]abci.Event)
	mapMsgWithdrawCommEvt := make(map[string]abci.Event)
	logger := ctx.Logger()
	listWallets := make(map[string]bool)
	slice := make([]string, 0)
	msgs := make([]*entity.Message, 0)
	msgSends := make([]*entity.MessageSend, 0)
	objTxs := make([]*entity.Txs, 0)
	blockResults, err := core.BlockResults(nil, &block.Height)
	if err != nil {
		logger.Error(err)
		return listWallets, nil, err
	}
	//counter := 0
	for idx, tx := range block.Txs {
		txResult := blockResults.TxsResults[idx]
		payloadLog := ""
		if !txResult.IsOK() {
			payloadLog = txResult.Log
		}
		sdkTx, err := txDecoder(tx)
		if err != nil {
			logger.Error(err)
			return listWallets, nil, err
		}
		stdTx, ok := sdkTx.(auth.StdTx)
		if !ok {
			logger.Error(err)
			return listWallets, nil, errors.New("cannot parse transaction")
		}
		// extract events in tx
		valoper := ""
		for _, e := range txResult.Events {
			idxReward := 0
			idxComm := 0
			switch e.Type {
			case "withdraw_rewards":
				for _, a := range e.Attributes {
					if string(a.Key) == "validator" {
						valoper = string(a.Value)
						keyReward := fmt.Sprintf("%v_%v", valoper, idxReward)
						mapMsgWithdrawRwdEvt[keyReward] = e
					}
				}
				idxReward += 1
			case "withdraw_commission":
				for _, a := range e.Attributes {
					if string(a.Key) == "amount" {
						keyComm := fmt.Sprintf("%v_%v", valoper, idxComm)
						mapMsgWithdrawCommEvt[keyComm] = e
					}
				}
				idxComm += 1
			}
		}

		// handle all messages in tx
		for i, m := range stdTx.Msgs {
			jsonMsg := cdc.MustMarshalJSON(m)
			objMsg := &entity.Message{
				BlockHeight:     block.Height,
				TransactionHash: hex.EncodeToString(tx.Hash()),
				MessageIndex:    int64(i),
				MessageTime:     block.Time.Unix(),
				Type:            m.Type(),
				IsValid:         txResult.IsOK(),
				Payload:         string(jsonMsg),
				PayloadLog:      payloadLog,
			}
			msgs = append(msgs, objMsg)
			if objMsg.IsValid {
				send, temp, err := ctrl.parseWalletsFromMessage(ctx, idx, m, objMsg)
				if err != nil {
					return listWallets, nil, err
				}
				if send != nil {
					msgSends = append(msgSends, send)
				}
				slice = appendWalletsToList(slice, temp...)
			}
		}
		txs, err := ctrl.preprocessTxs(ctx, block, idx, tx, txResult, stdTx, payloadLog, mapMsgWithdrawRwdEvt, mapMsgWithdrawCommEvt)
		if err != nil {
			logger.Error(err)
			return listWallets, nil, err
		}
		objTxs = append(objTxs, txs...)
	}
	limit := 2000
	for i := 0; i < len(msgSends); i += limit {
		batch := msgSends[i:utils.Min2Int(i+limit, len(msgSends))]
		if err := ctrl.messageRepo.MultiRowsInsertMessageSend(ctx, batch); err != nil {
			logger.Error(err)
			return listWallets, nil, err
		}
	}
	for i := 0; i < len(objTxs); i += limit {
		batch := objTxs[i:utils.Min2Int(i+limit, len(objTxs))]
		if err := ctrl.txsRepo.MultiRowsInsertTxs(ctx, batch); err != nil {
			logger.Error(err)
			return listWallets, nil, err
		}
	}

	appendListToMap(listWallets, slice)
	return listWallets, objTxs, nil
}

func (ctrl *Controller) parseWalletsFromMessage(ctx context.Context, txIndex int, msg cmtypes.Msg, objMsg *entity.Message) (*entity.MessageSend, []string, error) {
	listWallets := make([]string, 0)
	//a := new(entity.Message)
	var objMsgSend *entity.MessageSend
	switch m := msg.(type) {
	case bank.MsgSend:
		objMsgSend = &entity.MessageSend{
			BlockHeight:      objMsg.BlockHeight,
			TransactionHash:  objMsg.TransactionHash,
			TransactionIndex: int64(txIndex),
			MessageIndex:     objMsg.MessageIndex,
			MessageTime:      objMsg.MessageTime,
			Type:             objMsg.Type,
			IsValid:          objMsg.IsValid,
			Payload:          objMsg.Payload,
			PayloadLog:       objMsg.PayloadLog,
			FromAddress:      m.FromAddress.String(),
			ToAddress:        m.ToAddress.String(),
			Amount:           m.Amount.String(),
		}
		listWallets = appendWalletsToList(listWallets, m.FromAddress.String(), m.ToAddress.String())
	case bank.MsgMultiSend:
		for _, i := range m.Inputs {
			listWallets = appendWalletsToList(listWallets, i.Address.String())
		}
		for _, o := range m.Outputs {
			listWallets = appendWalletsToList(listWallets, o.Address.String())
		}
	case distribution.MsgWithdrawDelegatorReward:
		delegator, err := ctrl.delegatorRepo.GetDelegatorByDelegatorAndValidator(ctx, m.DelegatorAddress.String(), m.ValidatorAddress.String())
		if err != nil {
			return nil, listWallets, err
		}
		listWallets = appendWalletsToList(listWallets, m.DelegatorAddress.String(), m.ValidatorAddress.String())
		if delegator.WithdrawAddress != "" {
			listWallets = appendWalletsToList(listWallets, delegator.WithdrawAddress)
		}
	case distribution.MsgWithdrawValidatorCommission:
		delegators, err := ctrl.delegatorRepo.GetDelegatorByValidator(ctx, m.ValidatorAddress.String())
		if err != nil {
			return nil, listWallets, err
		}
		for _, i := range delegators {
			if i.WithdrawAddress != "" {
				listWallets = append(listWallets, i.WithdrawAddress)
			}
		}
		listWallets = appendWalletsToList(listWallets, m.ValidatorAddress.String())
	case distribution.MsgSetWithdrawAddress:
		obj := &entity.Delegator{
			DelegatorAddress: m.DelegatorAddress.String(),
			WithdrawAddress:  m.WithdrawAddress.String(),
		}
		listWallets = append(listWallets, m.WithdrawAddress.String(), m.DelegatorAddress.String())
		if err := ctrl.delegatorRepo.UpdateDelegatorWithdrawAddress(ctx, obj); err != nil {
			return nil, listWallets, err
		}
	case staking.MsgCreateValidator:
		listWallets = appendWalletsToList(listWallets, m.DelegatorAddress.String(), m.ValidatorAddress.String())
		validator, found := stakingKeeper.GetValidator(ctx.CMCtx(), m.ValidatorAddress)
		if !found {
			break
		}
		consPubkey, _ := cmtypes.Bech32ifyPubKey(cmtypes.Bech32PubKeyTypeConsPub, validator.ConsPubKey)
		val := &entity.Validator{
			Address:       m.ValidatorAddress.String(),
			Tokens:        m.Value.String(),
			Power:         validator.GetConsensusPower(),
			Jailed:        validator.Jailed,
			Status:        validator.Status.String(),
			ValConsAddr:   validator.GetConsAddr().String(),
			ValConsPubkey: consPubkey,
		}
		if err := ctrl.validatorRepo.CreateValidator(ctx, val); err != nil {
			return nil, listWallets, err
		}
		del := &entity.Delegator{
			DelegatorAddress: m.DelegatorAddress.String(),
			ValidatorAddress: m.ValidatorAddress.String(),
			Shares:           validator.DelegatorShares.String(),
		}
		if err := ctrl.delegatorRepo.CreateDelegator(ctx, del); err != nil {
			return nil, listWallets, err
		}
		decStake, _ := decimal.NewFromString(m.Value.Amount.String())
		stakeObj := &entity.Stake{
			DelegatorAddress: m.DelegatorAddress.String(),
			Shares:           decStake.String(),
			SharesDec:        decStake,
		}
		if err := ctrl.stakeRepo.UpsertStakedCoins(ctx, stakeObj); err != nil {
			return nil, listWallets, err
		}
	case staking.MsgEditValidator:
		obj := &entity.Validator{
			Address:  m.ValidatorAddress.String(),
			Detail:   m.Description.Details,
			Identity: m.Description.Identity,
			Moniker:  m.Description.Moniker,
			Website:  m.Description.Website,
		}
		if err := ctrl.validatorRepo.UpdateValidator(ctx, obj); err != nil {
			return nil, listWallets, err
		}
	case staking.MsgDelegate:
		listWallets = appendWalletsToList(listWallets, m.DelegatorAddress.String(), m.ValidatorAddress.String())
		if err := ctrl.updateStakingDelegate(ctx, m); err != nil {
			return nil, listWallets, err
		}
	case staking.MsgUndelegate:
		listWallets = appendWalletsToList(listWallets, m.DelegatorAddress.String(), m.ValidatorAddress.String())
		if err := ctrl.updateStakingUndelegate(ctx, m); err != nil {
			return nil, listWallets, err
		}
	case staking.MsgBeginRedelegate:
		listWallets = appendWalletsToList(listWallets, m.DelegatorAddress.String(), m.ValidatorSrcAddress.String(), m.ValidatorDstAddress.String())
		if err := ctrl.updateBeginRedelegate(ctx, m); err != nil {
			return nil, listWallets, err
		}
	case slashing.MsgUnjail:
		// update validator status in table validator
		if err := ctrl.updateMsgUnjail(ctx, m); err != nil {
			return nil, listWallets, err
		}
	}
	return objMsgSend, listWallets, nil
}

func appendWalletsToMapWallet(list map[string]bool, addresses ...string) map[string]bool {
	for _, a := range addresses {
		if !list[a] {
			list[a] = true
		}
	}
	return list
}

func appendWalletsToList(list []string, addresses ...string) []string {
	list = append(list, addresses...)
	return list
}

func appendListToMap(result map[string]bool, list []string) map[string]bool {
	for _, a := range list {
		if !result[a] {
			result[a] = true
		}
	}
	return result
}
