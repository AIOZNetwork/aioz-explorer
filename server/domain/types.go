package domain

const (
	MinPassLength          = 8
	MnemonicEntropySize    = 128
	DefaultBIP39Passphrase = ""

	AccAddressFormat       = "aioz"
	ValidatorAddressFormat = "aiozvaloper"

	Table_block                                 = "blocks"
	Table_delegator                             = "delegators"
	Table_message_begin_delegate                = "message_begin_delegates"
	Table_message_create_validator              = "message_create_validators"
	Table_message_delegate                      = "message_delegates"
	Table_message_multi_send                    = "message_multi_sends"
	Table_message_send                          = "message_sends"
	Table_message_undelegate                    = "message_undelegates"
	Table_message_withdraw_delegator_reward     = "message_withdraw_delegator_rewards"
	Table_message_withdraw_validator_commission = "message_withdraw_validator_commissions"
	Table_message                               = "messages"
	Table_pn_token_device                       = "pn_token_devices"
	Table_transaction                           = "transactions"
	Table_txs                                   = "txs"
	Table_validator                             = "validators"
	Table_wallet_address                        = "wallet_addresses"
	Table_stake                                 = "stakes"
)
