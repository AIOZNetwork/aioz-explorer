package ws

import (
	ws "aioz.io/go-aioz/cmd/aiozmedia/io/websocket"
	"encoding/json"
	"github.com/pkg/errors"
	"math/rand"
	"time"
)

type sendMessage struct {
	Namespace string `json:"namespace"`
	ClientID  int64  `json:"client_id"`
	Id        int64  `json:"id"`
	MsgID     int64  `json:"msg_id"`
	Result    interface{} `json:"result"`
}

func NewSendMessage(namespace string, id, clientID int64, typ string, data interface{}) (*sendMessage, error) {
	rand.Seed(time.Now().UnixNano())
	mid := rand.Int63n(999999999999)

	var body = struct {
		CID     int64       `json:"cid"`
		MsgType string      `json:"msg_type"`
		MsgID   int64       `json:"msg_id"`
		MsgData interface{} `json:"msg_data"`
	}{
		clientID,
		typ,
		mid,
		data,
	}

	return &sendMessage{
		Namespace: namespace,
		ClientID:  clientID,
		Id:        id,
		MsgID:     mid,
		Result:    body,
	}, nil
}

type receiveMessage struct {
	Namespace string
	ClientID  int64
	MsgData   interface{} `json:"msg_data"`
	MsgType   string      `json:"msg_type"`
	MsgID     int64       `json:"msg_id"`
}

func NewReceiveMessage(namespace string, clientID int64, body []byte) (*receiveMessage, error) {
	var msg *receiveMessage

	if err := json.Unmarshal(body, &msg); err != nil {
		return nil, errors.WithStack(err)
	}

	m := &receiveMessage{
		Namespace: namespace,
		ClientID:  clientID,
		MsgID:     msg.MsgID,
		MsgData:   msg.MsgData,
		MsgType:   msg.MsgType,
	}
	return m, nil
}

type wsRcvMessage struct {
	Namespace string
	ClientID  int64
	Id        uint32        `json:"id"`
	Msg       *ws.WsMessage `json:"msg"`
}

func NewWsRcvMessage(namespace string, clientID int64, body []byte) (*wsRcvMessage, error) {
	var msg *ws.WsMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		return nil, err
	}

	m := &wsRcvMessage{
		Namespace: namespace,
		ClientID:  clientID,
		Id:        0,
		Msg:       msg,
	}
	return m, nil
}
