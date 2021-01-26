package types

import (
	cmtypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

type GSAccount struct {
	Address        string        `json:"address"`
	Coins          cmtypes.Coins `json:"coins"`
	AccountNumber  string        `json:"account_number"`
	SequenceNumber string        `json:"sequence_number"`
}

type GenUtil struct {
	Gentxs []GenTxs `json:"gentxs"`
}

type GenTxs struct {
	Type  string `json:"type"`
	Value MsgObj `json:"value"`
}

type MsgObj struct {
	Message []Msg `json:"msg"`
}

type Msg struct {
	Type  string                     `json:"type"`
	Value staking.MsgCreateValidator `json:"value"`
}
