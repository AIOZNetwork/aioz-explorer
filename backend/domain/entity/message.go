package entity

type Message struct {
	BlockHeight     int64  `gorm:"primary_key;index:idx_message_block_height"`
	TransactionHash string `gorm:"primary_key"`
	MessageIndex    int64  `gorm:"primary_key"`
	MessageTime     int64  `gorm:"index:idx_message_message_time"`
	Type            string
	IsValid         bool
	Payload         string
	PayloadLog      string
}

func (m *Message) NewMessage(msg *Message) *Message {
	return msg
}

type MessageSend struct {
	BlockHeight      int64  `gorm:"primary_key;index:idx_message_block_height;index:idx_message_comp_msg,priority:1"`
	TransactionHash  string `gorm:"index:idx_message_transaction_hash"`
	TransactionIndex int64  `gorm:"primary_key;index:idx_message_transaction_index;index:idx_message_comp_msg,priority:2"`
	MessageIndex     int64  `gorm:"primary_key;index:idx_message_message_index;index:idx_message_comp_msg,priority:3"`
	MessageTime      int64  `gorm:"index:idx_message_message_time"`
	Type             string
	IsValid          bool
	Payload          string
	PayloadLog       string
	FromAddress      string `gorm:"index:idx_message_send_to_address"`
	ToAddress        string `gorm:"index:idx_message_send_from_address"`
	Amount           string
}

type MessageMultiSend struct {
	*Message
	Inputs  string
	Outputs string
}

type MessageDelegate struct {
	*Message
	DelegatorAddress string `gorm:"index:idx_message_delegate_delegator_address"`
	ValidatorAddress string `gorm:"index:idx_message_delegate_validator_address"`
	Amount           string
}

type MessageBeginDelegate struct {
	*Message
	DelegatorAddress    string `gorm:"index:idx_message_begin_delegate_delegator_address"`
	ValidatorSrcAddress string `gorm:"index:idx_message_begin_delegate_validator_src_address"`
	ValidatorDstAddress string `gorm:"index:idx_message_begin_delegate_validator_dst_address"`
	Amount              string
}

type MessageUndelegate struct {
	*Message
	DelegatorAddress string `gorm:"index:idx_message_undelegate_delegator_address"`
	ValidatorAddress string `gorm:"index:idx_message_undelegate_validator_address"`
	Amount           string
}

type MessageCreateValidator struct {
	*Message
	DelegatorAddress string `gorm:"index:idx_message_create_validator_delegator_address"`
	ValidatorAddress string `gorm:"index:idx_message_create_validator_validator_address"`
	Value            string
	Commission       string
}

type MessageWithdrawDelegatorReward struct {
	*Message
	DelegatorAddress string `gorm:"index:idx_message_withdraw_delegator_reward_delegator_address"`
	ValidatorAddress string `gorm:"index:idx_message_withdraw_delegator_reward_validator_address"`
}

type MessageWithdrawValidatorCommission struct {
	*Message
	ValidatorAddress string `gorm:"index:idx_message_withdraw_validator_commission_validator_address"`
}
