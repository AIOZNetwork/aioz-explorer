package ws

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type WS struct {
	Hub *hub
}

func NewWS() *WS {
	GlobalHub = newHub()
	return &WS{
		Hub: GlobalHub,
	}
}

func (n *WS) ServeWebsocket(c echo.Context) error {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if _, ok := err.(websocket.HandshakeError); ok {
		log.Println("ws: Not a websocket handshake")
		return c.JSON(http.StatusInternalServerError, "Invalid websocket handshake")
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	rand.Seed(time.Now().UnixNano())
	id := int64(rand.Intn(9999))
	n.Register("client", id, conn)

	return c.JSON(http.StatusOK, "websocket connected")
}

func (n *WS) Run() {
	n.Hub.run()
}

func (n *WS) Register(namespace string, id int64, conn *websocket.Conn) {
	client := newClient(namespace, id, n.Hub, conn)
	n.Hub.register <- client
	go client.writePump()
	client.readPump()
}

func (n *WS) IsConnecting(namespace string, clientID int64) bool {
	for c, _ := range n.Hub.Clients {
		if c.Namespace == namespace && c.Id == clientID {
			return true
		}
	}
	return false
}

func (n *WS) Notify(ctx context.Context, namespace string, id, clientID, userID int64, typ string, data interface{}) (msgID int64, errR error) {
	sendMsg, err := NewSendMessage(namespace, id, clientID, userID, typ, data)
	if err != nil {
		return 0, err
	}

	n.Hub.Send <- sendMsg

	return sendMsg.MsgID, nil
}

func (n *WS) WaitResponse(ctx context.Context, msgID int64) (data interface{}, errR error) {
	return n.Hub.waitResponse(ctx, msgID)
}

func (n *WS) MustNotify(ctx context.Context, namespace string, id, clientID, userID int64, typ string, data interface{}) (interface{}, error) {
	msgID, err := n.Notify(ctx, namespace, id, clientID, userID, typ, data)
	if err != nil {
		return nil, err
	}
	return n.WaitResponse(ctx, msgID)
}

//func (n *WS) Broadcast(namespace string, typ string, data interface{}) {
//	clients := n.Hub.findClientsOfNamespace(namespace)
//	for _, c := range clients {
//		_, _ = n.Notify(context.Background(), c.Namespace, c.Id, typ, data)
//	}
//}
