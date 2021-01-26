package controller

import (
	"aioz.io/go-aioz/types"
	"aioz.io/go-aioz/x_gob_explorer/context"
	"aioz.io/go-aioz/x_gob_explorer/domain/entity"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

func (ctrl *Controller) preprocessTxs(ctx context.Context, block *tmtypes.Block, txIndex int, tx tmtypes.Tx, txResult *abci.ResponseDeliverTx,
	stdTx auth.StdTx, payloagLog string, mapReward, mapComm map[string]abci.Event) ([]*entity.Txs, error) {
	result := make([]entity.Txs, 0)
	idxReward := 0
	idxComm := 0
	withdrawAddressReward := ""
	for i, msg := range stdTx.Msgs {
		jsonMsg := cdc.MustMarshalJSON(msg)
		obj := entity.Txs{
			BlockHeight:      block.Height,
			BlockHash:        block.AppHash.String(),
			BlockTime:        block.Time.Unix(),
			TransactionHash:  hex.EncodeToString(tx.Hash()),
			TransactionIndex: int64(txIndex),
			MessageType:      msg.Type(),
			MessageIndex:     int64(i),
			Address:          "",
			AddressFrom:      make([]string, 0),
			AddressTo:        make([]string, 0),
			Amount:           "",
			IsValid:          txResult.IsOK(),
			PayloadMsg:       string(jsonMsg),
			PayloadErr:       payloagLog,
		}
		switch m := msg.(type) {
		case bank.MsgSend:
			obj.AddressFrom = append(obj.AddressFrom, m.FromAddress.String())
			obj.AddressTo = append(obj.AddressTo, m.ToAddress.String())
			amt, _ := json.Marshal(m.Amount)
			obj.Amount = string(amt)
			obj.Address = m.FromAddress.String()
			result = append(result, obj)
			obj.Address = m.ToAddress.String()
			result = append(result, obj)
		case bank.MsgMultiSend:
			for _, f := range m.Inputs {
				obj.AddressFrom = append(obj.AddressFrom, f.Address.String())
				obj.Address = f.Address.String()
				amt, _ := json.Marshal(f.Coins)
				obj.Amount = string(amt)
				result = append(result, obj)
			}
			for _, f := range m.Outputs {
				obj.AddressTo = append(obj.AddressTo, f.Address.String())
				obj.Address = f.Address.String()
				amt, _ := json.Marshal(f.Coins)
				obj.Amount = string(amt)
				result = append(result, obj)
			}
		case staking.MsgCreateValidator:
			obj.AddressFrom = append(obj.AddressFrom, m.DelegatorAddress.String())
			obj.AddressTo = append(obj.AddressTo, m.ValidatorAddress.String())
			amt, _ := json.Marshal(m.Value)
			obj.Amount = string(amt)
			obj.Address = m.DelegatorAddress.String()
			result = append(result, obj)
			obj.Address = m.ValidatorAddress.String()
			result = append(result, obj)
		case distribution.MsgSetWithdrawAddress:
			obj.Address = m.DelegatorAddress.String()
			obj.AddressFrom = append(obj.AddressFrom, m.DelegatorAddress.String())
			obj.AddressTo = append(obj.AddressTo, m.WithdrawAddress.String())
			result = append(result, obj)
			obj.Address = m.WithdrawAddress.String()
			result = append(result, obj)
		case staking.MsgDelegate:
			obj.AddressFrom = append(obj.AddressFrom, m.DelegatorAddress.String())
			obj.AddressTo = append(obj.AddressTo, m.ValidatorAddress.String())
			amt, _ := json.Marshal(m.Amount)
			obj.Amount = string(amt)
			obj.Address = m.DelegatorAddress.String()
			result = append(result, obj)
			obj.Address = m.ValidatorAddress.String()
			result = append(result, obj)
		case staking.MsgUndelegate:
			obj.Address = m.DelegatorAddress.String()
			obj.AddressFrom = append(obj.AddressFrom, m.DelegatorAddress.String())
			obj.AddressTo = append(obj.AddressTo, m.ValidatorAddress.String())
			amt, _ := json.Marshal(m.Amount)
			obj.Amount = string(amt)
			result = append(result, obj)
			obj.Address = m.ValidatorAddress.String()
			result = append(result, obj)
		case staking.MsgBeginRedelegate:
			obj.AddressFrom = append(obj.AddressFrom, m.DelegatorAddress.String())
			obj.AddressTo = append(obj.AddressTo, m.ValidatorSrcAddress.String(), m.ValidatorDstAddress.String())
			obj.Address = m.DelegatorAddress.String()
			amt, _ := json.Marshal(m.Amount)
			obj.Amount = string(amt)
			result = append(result, obj)
			obj.Address = m.ValidatorDstAddress.String()
			result = append(result, obj)
			obj.Address = m.ValidatorSrcAddress.String()
			result = append(result, obj)
		case distribution.MsgWithdrawDelegatorReward:
			keyReward := fmt.Sprintf("%v_%v", m.ValidatorAddress.String(), idxReward)
			if v, ok := mapReward[keyReward]; ok {
				for _, attr := range v.Attributes {
					if string(attr.Key) == "amount" {
						amount, _ := types.ParseCoins(string(attr.Value))
						amt, _ := json.Marshal(amount)
						obj.Amount = string(amt)
					}
				}
			}
			delegator, err := ctrl.delegatorRepo.GetDelegatorByDelegatorAndValidator(ctx, m.DelegatorAddress.String(), m.ValidatorAddress.String())
			if err != nil {
				return nil, err
			}
			if delegator.WithdrawAddress != "" {
				obj.AddressTo = append(obj.AddressTo, delegator.WithdrawAddress)
				withdrawAddressReward = delegator.WithdrawAddress
			} else {
				obj.AddressTo = append(obj.AddressTo, m.DelegatorAddress.String())
				withdrawAddressReward = m.DelegatorAddress.String()
			}
			obj.Address = m.DelegatorAddress.String()
			result = append(result, obj)
			obj.Address = m.ValidatorAddress.String()
			result = append(result, obj)
			obj.AddressFrom = append(obj.AddressFrom, m.DelegatorAddress.String())
			idxReward += 1
		case distribution.MsgWithdrawValidatorCommission:
			keyComm := fmt.Sprintf("%v_%v", m.ValidatorAddress.String(), idxComm)
			if v, ok := mapComm[keyComm]; ok {
				for _, attr := range v.Attributes {
					if string(attr.Key) == "amount" {
						amount, _ := types.ParseCoins(string(attr.Value))
						amt, _ := json.Marshal(amount)
						obj.Amount = string(amt)
					}
				}
			}
			obj.Address = m.ValidatorAddress.String()
			obj.AddressFrom = append(obj.AddressFrom, m.ValidatorAddress.String())
			obj.AddressTo = append(obj.AddressTo, withdrawAddressReward)
			result = append(result, obj)
			idxComm += 1
		case slashing.MsgUnjail:
			obj.AddressFrom = append(obj.AddressFrom, m.ValidatorAddr.String())
			obj.Address = m.ValidatorAddr.String()
			result = append(result, obj)
		case staking.MsgEditValidator:
			from := ""
			if len(m.GetSigners()) > 0 {
				from = m.GetSigners()[0].String()
			}
			if from != "" {
				obj.Address = from
				obj.AddressFrom = append(obj.AddressFrom, from)
				obj.AddressTo = append(obj.AddressFrom, m.ValidatorAddress.String())
				result = append(result, obj)
				obj.Address = m.ValidatorAddress.String()
				result = append(result, obj)
			} else {
				obj.Address = m.ValidatorAddress.String()
				obj.AddressFrom = append(obj.AddressFrom, m.ValidatorAddress.String())
				result = append(result, obj)
			}
		default:
			signer := stdTx.GetSigners()[0].String()
			obj.Address = signer
			obj.AddressFrom = append(obj.AddressFrom, signer)
			amt, _ := json.Marshal(stdTx.Fee.Amount)
			obj.Amount = string(amt)
			result = append(result, obj)
		}
	}

	a := make([]*entity.Txs, 0)
	for _, tx := range result {
		a = append(a, &entity.Txs{
			BlockHeight:      tx.BlockHeight,
			BlockHash:        tx.BlockHash,
			BlockTime:        tx.BlockTime,
			TransactionHash:  tx.TransactionHash,
			TransactionIndex: tx.TransactionIndex,
			MessageType:      tx.MessageType,
			MessageIndex:     tx.MessageIndex,
			Address:          tx.Address,
			AddressFrom:      tx.AddressFrom,
			AddressTo:        tx.AddressTo,
			Amount:           tx.Amount,
			IsValid:          tx.IsValid,
			PayloadMsg:       tx.PayloadMsg,
			PayloadErr:       tx.PayloadErr,
		})
	}
	return a, nil
}
