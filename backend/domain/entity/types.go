package entity

type MultiSendInputs []MultiSendInput
type MultiSendOutputs []MultiSendOutput

type MultiSendInput struct {
	Address string
	Amount  string
}

type MultiSendOutput struct {
	Address string
	Amount  string
}

type CommRate struct {
	Rate          string
	MaxRate       string
	MaxChangeRate string
}
