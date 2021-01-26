package blockchain

import (
	"aioz.io/go-aioz/x_gob_explorer/domain/entity"
	"aioz.io/go-aioz/x_gob_explorer/ws"
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func subscribeWallets(wshandler *ws.WS, msgs []*entity.Txs) {
	for _, txs := range msgs {
		if wshandler.Hub.Subscribers[txs.Address] {
			// if there is address to be concerned -> notify to client
			clients := wshandler.Hub.FindClientOfWallet(txs.Address)
			for _, client := range clients {
				var txtype string
				var coins sdk.Coins
				if sliceContains(txs.Address, txs.AddressFrom) {
					txtype = "send"
				} else if sliceContains(txs.Address, txs.AddressTo) {
					txtype = "receive"
				}
				err := json.Unmarshal([]byte(txs.Amount), &coins)
				if err != nil {
					coins = sdk.NewCoins()
				}
				resp := &entity.TxsResp{
					BlockHeight:      txs.BlockHeight,
					BlockHash:        txs.BlockHash,
					BlockTime:        txs.BlockTime,
					TransactionHash:  txs.TransactionHash,
					TransactionIndex: txs.TransactionIndex,
					TransactionType:  txtype,
					MessageType:      txs.MessageType,
					MessageIndex:     txs.MessageIndex,
					Address:          txs.Address,
					AddressFrom:      txs.AddressFrom,
					AddressTo:        txs.AddressTo,
					Amount:           coins,
					IsValid:          txs.IsValid,
					PayloadMsg:       txs.PayloadMsg,
					PayloadErr:       txs.PayloadErr,
				}
				for k, _ := range client.MapUserWallet {
					msg, _ := ws.NewSendMessage(client.Namespace, 0, client.Id, k, "wallet.response", resp)
					wshandler.Hub.Send <- msg
				}
			}
		}

	}
}

func subscribeMessageSend(wshandler *ws.WS, msgs []*entity.Txs) {
	for _, m := range msgs {
		if wshandler.Hub.Subscribers[m.MessageType] {
			// if there is address to be concerned -> notify to client
			clients := wshandler.Hub.FindClientOfMsgType(m.MessageType)
			for _, client := range clients {
				for k, _ := range client.MapUserMsg {
					msg, _ := ws.NewSendMessage(client.Namespace, 0, client.Id, k, "message.response", m)
					wshandler.Hub.Send <- msg
				}
			}
		}
	}
}

func subscribeBlock(wshandler *ws.WS, block interface{}) {
	for client, _ := range wshandler.Hub.Clients {
		for c, v := range client.MapUserBlock {
			if v {
				msg, _ := ws.NewSendMessage(client.Namespace, 0, client.Id, c, "block.response", block)
				wshandler.Hub.Send <- msg
			}
		}
	}
}

func sliceContains(input string, array []string) bool {
	for _, a := range array {
		if a == input {
			return true
		}
	}
	return false
}
