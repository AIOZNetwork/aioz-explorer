package message

import (
	mediatypes "aioz.io/go-aioz/cmd/aiozmedia/types"
	"encoding/json"
)

var (
	_ mediatypes.IMessage = &ClientSubscription{}
	_ mediatypes.IResult  = &ClientResponse{}
)

type ClientSubscription struct {
	CID     int64       `json:"cid"`
	MsgId   int64       `json:"msg_id"`
	MsgType string      `json:"msg_type"`
	MsgData interface{} `json:"msg_data"`
}

type ClientResponse struct {
	ClientId int64       `json:"client_id"`
	MsgId    int64       `json:"msg_id"`
	MsgType  string      `json:"msg_type"`
	MsgData  interface{} `json:"msg_data"`
	CID      int64       `json:"cid"`
}

func (c *ClientResponse) ToByte() []byte {
	ret, _ := json.Marshal(c)
	return ret
}
