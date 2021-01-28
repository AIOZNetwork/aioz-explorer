package domain

import sdk "github.com/cosmos/cosmos-sdk/types"

type Message struct {
	BlockHeight     int64  `gorm:"primary_key" json:"block_height"`
	TransactionHash string `gorm:"primary_key" json:"transaction_hash"`
	MessageIndex    int64  `gorm:"primary_key" json:"message_index"`
	MessageTime     uint64 `json:"message_time"`
	Type            string `json:"type"`
	IsValid         bool   `json:"is_valid"`
	Payload         string
	PayloadLog      string
}

func (m *Message) NewMessage(msg *Message) *Message {
	return msg
}

type MessageSend struct {
	BlockHeight     int64  `gorm:"primary_key" json:"block_height"`
	TransactionHash string `gorm:"primary_key" json:"transaction_hash"`
	MessageIndex    int64  `gorm:"primary_key" json:"message_index"`
	MessageTime     uint64 `json:"message_time"`
	Type            string `json:"type"`
	IsValid         bool   `json:"is_valid"`
	Payload         string
	PayloadLog      string
	TransactionType string `json:"transaction_type"`
	FromAddress     string `json:"from_address"`
	ToAddress       string `json:"to_address"`
	Amount          string `json:"amount"`
}

type MessageSendResp struct {
	*Message
	TransactionType string    `json:"transaction_type"`
	FromAddress     string    `json:"from_address"`
	ToAddress       string    `json:"to_address"`
	Amount          sdk.Coins `json:"amount"`
}

type MessageMultiSend struct {
	*Message
	Inputs  string
	Outputs string
}

type MessageDelegate struct {
	*Message
	DelegateAddress  string
	ValidatorAddress string
	Amount           string
}

type MessageBeginDelegate struct {
	*Message
	DelegateAddress     string
	ValidatorSrcAddress string
	ValidatorDstAddress string
	Amount              string
}

type MessageUndelegate struct {
	*Message
	DelegateAddress  string
	ValidatorAddress string
	Amount           string
}

type MessageCreateValidator struct {
	*Message
	DelegatorAddress string
	ValidatorAddress string
	Value            string
	Commission       string
}
