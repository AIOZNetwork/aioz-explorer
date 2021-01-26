package ws

import (
	"encoding/json"
	"github.com/pkg/errors"
	"math/rand"
	"time"
)

type sendMessage struct {
	Namespace string      `json:"namespace"`
	Id        int64       `json:"id"`
	ClientId  int64       `json:"client_id"`
	MsgID     int64       `json:"msg_id"`
	Result    interface{} `json:"result"`
}

func NewSendMessage(namespace string, id, clientID, userSubscriptionID int64, typ string, data interface{}) (*sendMessage, error) {
	rand.Seed(time.Now().UnixNano())
	mid := rand.Int63n(999999999999)

	var body = struct {
		CID     int64       `json:"cid"`
		MsgType string      `json:"msg_type"`
		MsgID   int64       `json:"msg_id"`
		MsgData interface{} `json:"msg_data"`
	}{
		userSubscriptionID,
		typ,
		mid,
		data,
	}

	//b, err := json.Marshal(body)
	//if err != nil {
	//	return nil, errors.WithStack(err)
	//}

	return &sendMessage{
		Namespace: namespace,
		Id:        id,
		ClientId:  clientID,
		MsgID:     mid,
		Result:    body,
	}, nil
}

type rcvWrapper struct {
	Id  int64          `json:"id"`
	Msg receiveMessage `json:"msg"`
}

type receiveMessage struct {
	Id        int64 `json:"id"`
	Namespace string
	ClientID  int64
	MsgData   interface{} `json:"msg_data"`
	MsgType   string      `json:"msg_type"`
	MsgID     int64       `json:"msg_id"`
	CID       int64       `json:"cid"`
}

func NewReceiveMessage(namespace string, clientID int64, body []byte) (*receiveMessage, error) {
	var msg *rcvWrapper

	if err := json.Unmarshal(body, &msg); err != nil {
		return nil, errors.WithStack(err)
	}

	m := &receiveMessage{
		Id:        msg.Id,
		Namespace: namespace,
		ClientID:  clientID,
		MsgID:     msg.Msg.MsgID,
		MsgData:   msg.Msg.MsgData,
		MsgType:   msg.Msg.MsgType,
		CID:       msg.Msg.CID,
	}
	return m, nil
}
