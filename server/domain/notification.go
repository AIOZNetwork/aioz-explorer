package domain

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/viper"
	"swagger-server/utils"
)

const (
	MsgSend = "send"
)

type Platform int

const (
	Platform_p_none    Platform = 0
	Platform_p_ios     Platform = 1
	Platform_p_android Platform = 2
	Platform_p_web     Platform = 11
	Platform_p_windows Platform = 12
)

var Platform_name = map[Platform]string{
	0:  "p_none",
	1:  "p_ios",
	2:  "p_android",
	11: "p_web",
	12: "p_windows",
}

var Platform_value = map[string]Platform{
	"p_none":    0,
	"p_ios":     1,
	"p_android": 2,
	"p_web":     11,
	"p_windows": 12,
}

type DeviceStatus int

const (
	Device_status_deactive = 0
	Device_status_active   = 1

	MessageNotiReceive = "recipient"
	MessageNotiSend    = "sender"

	DefaultDecCoin = "1000000000000000000"
)

type NotificationPayload interface {
	Jsonify() string
	GetType() string
	GetAddress() string
	GetCoins() string
}

type NotificationMsg struct {
	BlockHeight     int64     `json:"block_height"`
	TxHash          string    `json:"transaction_hash"`
	TxType          string    `json:"type"`
	TransactionType string    `json:"transaction_type"`
	FromAddress     []string  `json:"from_address"`
	ToAddress       []string  `json:"to_address"`
	Amount          sdk.Coins `json:"amount"`
	MessageIndex    int64     `json:"message_index"`
	MessageTime     int64     `json:"message_time"`
}

func NewNotificationMsg(t TxsResp) NotificationMsg {
	return NotificationMsg{
		BlockHeight:     t.BlockHeight,
		TxHash:          t.TransactionHash,
		TxType:          "",
		TransactionType: t.TransactionType,
		FromAddress:     t.AddressFrom,
		ToAddress:       t.AddressTo,
		Amount:          t.Amount,
		MessageIndex:    t.MessageIndex,
		MessageTime:     t.BlockTime,
	}
}

func (n NotificationMsg) Jsonify() string {
	b, _ := json.Marshal(n)
	return string(b)
}

func (n NotificationMsg) GetType() string {
	switch n.TransactionType {
	case "send":
		return MessageNotiSend
	case "receive":
		return MessageNotiReceive
	default:
		return ""
	}
}

func (n NotificationMsg) GetAddress() string {
	switch n.TransactionType {
	case "receive":
		if len(n.FromAddress) > 0 {
			return n.FromAddress[0]
		}
		return ""
	case "send":
		if len(n.ToAddress) > 0 {
			return n.ToAddress[0]
		}
		return ""
	default:
		return ""
	}
}

func (n NotificationMsg) GetCoins() string {
	coins := sdk.NewDecCoins(n.Amount)
	denom := viper.GetString("coins.denom")
	token := viper.GetString("coins.token")
	for i, _ := range n.Amount {
		if coins[i].Denom == denom {
			amt, err := sdk.NewDecFromStr(DefaultDecCoin)
			if err != nil {
				continue
			}
			coins[i].Amount = coins[i].Amount.Quo(amt)
			coins[i].Denom = token
		}
	}
	return utils.RemoveTrailingZerosFromDecCoins(coins)
}

type NotificationMessageRcv struct {
	BlockHeight     int64     `json:"block_height"`
	TxHash          string    `json:"transaction_hash"`
	TxType          string    `json:"type"`
	TransactionType string    `json:"transaction_type"`
	FromAddress     string    `json:"from_address"`
	ToAddress       string    `json:"to_address"`
	Amount          sdk.Coins `json:"amount"`
	MessageIndex    int64     `json:"message_index"`
	MessageTime     int64     `json:"message_time"`
}

func (nm NotificationMessageRcv) Jsonify() string {
	b, err := json.Marshal(nm)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (nm NotificationMessageRcv) GetType() string {
	return MessageNotiReceive
}

func (nm NotificationMessageRcv) GetAddress() string {
	return nm.FromAddress
}

func (nm NotificationMessageRcv) GetCoins() string {
	coins := sdk.NewDecCoins(nm.Amount)

	for i, _ := range nm.Amount {
		amt, err := sdk.NewDecFromStr(DefaultDecCoin)
		if err != nil {
			continue
		}
		coins[i].Amount = coins[i].Amount.Quo(amt)
	}
	return utils.RemoveTrailingZerosFromDecCoins(coins)
}

type NotificationMessageSnd struct {
	BlockHeight     int64     `json:"block_height"`
	TxHash          string    `json:"transaction_hash"`
	TxType          string    `json:"type"`
	TransactionType string    `json:"transaction_type"`
	FromAddress     string    `json:"from_address"`
	ToAddress       string    `json:"to_address"`
	Amount          sdk.Coins `json:"amount"`
	MessageIndex    int64     `json:"message_index"`
	MessageTime     int64     `json:"message_time"`
}

func (nm NotificationMessageSnd) Jsonify() string {
	b, err := json.Marshal(nm)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (nm NotificationMessageSnd) GetType() string {
	return MessageNotiSend
}

func (nm NotificationMessageSnd) GetAddress() string {
	return nm.ToAddress
}

func (nm NotificationMessageSnd) GetCoins() string {
	coins := sdk.NewDecCoins(nm.Amount)

	for i, _ := range nm.Amount {
		amt, err := sdk.NewDecFromStr(DefaultDecCoin)
		if err != nil {
			continue
		}
		coins[i].Amount = coins[i].Amount.QuoTruncate(amt)
	}
	return utils.RemoveTrailingZerosFromDecCoins(coins)
}
