package ws

import (
	context2 "context"
	firebase "firebase.google.com/go/v4"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"log"
	"math/rand"
	"net/http"
	"swagger-server/context"
	"swagger-server/ws/message"
	"time"
)

type WS struct {
	hub *hub
}

func NewWS(ctx context.Context, wsc *WSClient, fbapp *firebase.App) *WS {
	GlobalHub = newHub(ctx, wsc, fbapp)
	return &WS{
		hub: GlobalHub,
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

func (n *WS) ResubscribeDeviceToken(wallets []string) {
	rand.Seed(time.Now().UnixNano())
	mid := rand.Int63n(999999999999)
	for {
		_, err := n.hub.WsClient.Send(context2.Background(), message.ClientSubscription{
			CID:     -5,
			MsgId:   mid,
			MsgType: "wallet.subscribe",
			MsgData: wallets,
		})
		if err != nil {
			log.Println("Send subscription wallet device token failed. Retrying...")
			n.hub.ctx.Logger().Error("Send subscription wallet device token failed. Retrying...")
			continue
		} else {
			log.Println("Re-subscribe wallet for devices successfully")
			n.hub.ctx.Logger().Info("Re-subscribe wallet for devices successfully")
			break
		}
	}
}

func (n *WS) Run() {
	n.hub.run()
}

func (n *WS) Register(namespace string, id int64, conn *websocket.Conn) {
	client := newClient(namespace, id, n.hub, conn)
	n.hub.register <- client
	go client.writePump()
	client.readPump()
}

func (n *WS) IsConnecting(namespace string, clientID int64) bool {
	for c, _ := range n.hub.clients {
		if c.namespace == namespace && c.id == clientID {
			return true
		}
	}
	return false
}

func (n *WS) Notify(ctx context2.Context, namespace string, id, clientID int64, typ string, data interface{}) (msgID int64, errR error) {
	sendMsg, err := NewSendMessage(namespace, id, clientID, typ, data)
	if err != nil {
		return 0, err
	}

	n.hub.send <- sendMsg

	return sendMsg.MsgID, nil
}

func (n *WS) WaitResponse(ctx context.Context, msgID int64) (data interface{}, errR error) {
	return n.hub.waitResponse(ctx, msgID)
}

func (n *WS) MustNotify(ctx context.Context, namespace string, id, clientID int64, typ string, data interface{}) (interface{}, error) {
	msgID, err := n.Notify(ctx, namespace, id, clientID, typ, data)
	if err != nil {
		return nil, err
	}
	return n.WaitResponse(ctx, msgID)
}

func (n *WS) Broadcast(namespace string, typ string, data interface{}) {
	clients := n.hub.findClientsOfNamespace(namespace)
	for _, c := range clients {
		_, _ = n.Notify(context2.Background(), c.namespace, 0, c.id, typ, data)
	}
}

func (n *WS) DisconnectAllClients() {
	for c, _ := range n.hub.clients {
		n.hub.unregister <- c
	}
}
