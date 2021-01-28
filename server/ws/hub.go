package ws

import (
	"encoding/json"
	firebase "firebase.google.com/go/v4"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"log"
	"math/rand"
	"swagger-server/context"
	"swagger-server/domain"
	firebasegob "swagger-server/firebase"
	"swagger-server/utils"
	"swagger-server/ws/message"
	"sync"
	"time"
)

var (
	GlobalHub *hub
)

type ack struct {
	mu sync.Mutex
	m  map[int64]chan interface{}
}

func newAck() *ack {
	return &ack{
		m: make(map[int64]chan interface{}),
	}
}

type hub struct {
	ctx context.Context

	mu          sync.Mutex
	WsClient    *WSClient
	FBApp       *firebase.App
	clients     map[*client]bool
	register    chan *client
	unregister  chan *client
	send        chan *sendMessage
	receive     chan *receiveMessage
	rcvWsClient chan *wsRcvMessage
	ack         *ack
}

func newHub(ctx context.Context, wsc *WSClient, fbapp *firebase.App) *hub {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	return &hub{
		ctx:         ctx,
		WsClient:    wsc,
		FBApp:       fbapp,
		clients:     make(map[*client]bool),
		register:    make(chan *client),
		unregister:  make(chan *client),
		send:        make(chan *sendMessage, 1000),
		receive:     make(chan *receiveMessage, 1000),
		rcvWsClient: make(chan *wsRcvMessage, 1000),
		ack:         newAck(),
	}
}

func (h *hub) waitResponse(ctx context.Context, msgID int64) (interface{}, error) {
	waiter := make(chan interface{})
	h.ack.m[msgID] = waiter
	defer func() {
		h.ack.mu.Lock()
		if _, ok := h.ack.m[msgID]; ok {
			delete(h.ack.m, msgID)
		}
		h.ack.mu.Unlock()
	}()

	select {
	case v := <-waiter:
		return v, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (h *hub) findClient(namespace string, id int64) *client {
	for c, _ := range h.clients {
		if c.id == id && c.namespace == namespace {
			return c
		}
	}
	return nil
}

func (h *hub) findClientsOfNamespace(namespace string) []*client {
	result := make([]*client, 0)
	for c, _ := range h.clients {
		if namespace == "all" || c.namespace == namespace {
			result = append(result, c)
		}
	}
	return result
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			h.mu.Lock()
			if old := h.findClient(c.namespace, c.id); old != nil {
				close(old.send)
				delete(h.clients, old)
			} else {
				fmt.Printf("Connect %d\n", c.id)
				h.ctx.Logger().Logf(logrus.InfoLevel, "[hub.connect] Client connected %v", c.id)
			}
			h.clients[c] = true
			h.mu.Unlock()

		case c := <-h.unregister:
			if _, ok := h.clients[c]; ok {
				h.mu.Lock()
				close(c.send)
				delete(h.clients, c)
				h.mu.Unlock()

				// unsubscribe to node websocket
				h.UnsubscribeWalletNodeWs(c)
				h.UnsubscribeMsgNodeWs(c)
				h.UnsubscribeBlockWs(c)

				fmt.Printf("Disconnect %d\n", c.id)
				h.ctx.Logger().Logf(logrus.InfoLevel, "[hub.disconnect] Client disconnected %v", c.id)
			}
		case msg := <-h.send:
			for c := range h.clients {
				if msg.Namespace == c.namespace && msg.ClientID == c.id {
					body, _ := json.Marshal(msg)
					select {
					case c.send <- body:
					default:
						h.mu.Lock()
						close(c.send)
						delete(h.clients, c)
						h.mu.Unlock()
					}
					h.ctx.Logger().Logf(logrus.InfoLevel, "[hub.send] Send message to client %v. Message content: %v", c.id, string(body))
				}
			}
		case msg := <-h.receive:
			client := h.findClient(msg.Namespace, msg.ClientID)
			jsonobj, _ := json.Marshal(msg.MsgData)
			data := make([]string, 0)
			_ = json.Unmarshal(jsonobj, &data)
			if msg.MsgType == "wallet.subscribe" {
				log.Println("Message ID: ", msg.MsgID)
				for _, d := range data {
					client.wallets = append(client.wallets, d)
				}
				h.SubscribeWalletNodeWs(client, 0)
			} else if msg.MsgType == "message.subscribe" {
				log.Println("Message ID: ", msg.MsgID)
				for _, d := range data {
					client.msgs = append(client.msgs, d)
				}
				h.SubscribeMsgNodeWs(client, 0)
			} else if msg.MsgType == "block.subscribe" {
				h.SubscribeBlockWs(client, 0)
			}
			h.ctx.Logger().Logf(logrus.InfoLevel, "[hub.receive] Server received subscription from client %v and forward to node. Message content: %v", client.id, msg)
		case msg := <-h.rcvWsClient:
			var r *message.ClientResponse
			if msg.Msg.Msg == nil {
				continue
			}
			b, _ := json.Marshal(msg.Msg.Msg)
			if err := json.Unmarshal(b, &r); err != nil {
				fmt.Println(err)
				continue
			}
			jsonobj, _ := json.Marshal(r.MsgData)
			data := make([]string, 0)
			_ = json.Unmarshal(jsonobj, &data)
			client := h.findClient(msg.Namespace, msg.ClientID)
			if r.MsgType == "wallet.subscribe" {
				for _, d := range data {
					client.wallets = append(client.wallets, d)
				}
				h.SubscribeWalletNodeWs(client, int64(msg.Msg.Id))
			} else if r.MsgType == "message.subscribe" {
				for _, d := range data {
					client.msgs = append(client.msgs, d)
				}
				h.SubscribeMsgNodeWs(client, int64(msg.Msg.Id))
			} else if r.MsgType == "block.subscribe" {
				h.SubscribeBlockWs(client, int64(msg.Msg.Id))
			}
			h.ctx.Logger().Logf(logrus.InfoLevel, "[hub.rcvWsClient] Server received subscription from client %v. Message content: %v", client.id, msg)
		case n := <-h.WsClient.NotificationCh:
			raw, _ := json.Marshal(n.Result)
			var res *message.ClientResponse
			err := json.Unmarshal(raw, &res)
			if err != nil {
				fmt.Println(err)
				continue
			}
			var r domain.TxsResp
			b, _ := json.Marshal(res.MsgData)
			_ = json.Unmarshal(b, &r)
			m, _ := NewSendMessage("client", 0, res.CID, res.MsgType, res.MsgData)
			go func() {
				h.send <- m
			}()
			h.ctx.Logger().Logf(logrus.InfoLevel, "[hub.NotificationCh] Server received response from node. Forward notification to client %v. Message content: %v", res.CID, m)
			if res.CID < 0 { // only send firebase notification to client which CID < 0
				switch res.MsgType {
				// send to ws client who subscribed before
				// besides, we should push notification to devices which subscribed before also
				// must cross-check wallet in MapWalletNotification vs response returned
				case "wallet.response":
					for k, v := range h.WsClient.MapWalletNotification {
						if utils.CheckExistsInSlice(v, r.Address) {
							m := domain.NewNotificationMsg(r)
							if err := firebasegob.PushNotifications(h.FBApp, []string{k}, m); err != nil {
								log.Println(err)
							}
							log.Println("Sending notification to device token: ", k)
							h.ctx.Logger().Logf(logrus.InfoLevel, "[hub.NotificationCh] Server received response from node. Forward notification to device token: %v. Message content: %v", k, m)
						}
					}
				case "message.response":
				default:
				}
			}
		}
	}
}

func (h *hub) SubscribeWalletNodeWs(c *client, rID int64) {
	rand.Seed(time.Now().UnixNano())
	mid := rand.Int63n(999999999999)
	result, err := h.WsClient.Send(h.ctx, message.ClientSubscription{
		CID:     c.id,
		MsgId:   mid,
		MsgType: "wallet.subscribe",
		MsgData: c.wallets,
	})
	if err != nil {
		log.Println(err)
	}
	if result != nil {
		raw, _ := json.Marshal(result)
		var res *message.ClientResponse
		json.Unmarshal(raw, &res)
		go func() {
			m, _ := NewSendMessage("client", rID, c.id, "ack.subscription", res.MsgData)
			h.send <- m
		}()
	}
}

func (h *hub) UnsubscribeWalletNodeWs(c *client) {
	rand.Seed(time.Now().UnixNano())
	mid := rand.Int63n(999999999999)
	_, err := h.WsClient.Send(h.ctx, message.ClientSubscription{
		CID:     c.id,
		MsgId:   mid,
		MsgType: "wallet.unsubscribe",
		MsgData: c.wallets,
	})
	if err != nil {
		log.Println(err)
	}
}

func (h *hub) SubscribeMsgNodeWs(c *client, rID int64) {
	rand.Seed(time.Now().UnixNano())
	mid := rand.Int63n(999999999999)
	result, err := h.WsClient.Send(h.ctx, message.ClientSubscription{
		CID:     c.id,
		MsgId:   mid,
		MsgType: "message.subscribe",
		MsgData: c.msgs,
	})
	if err != nil {
		log.Println(err)
	}
	if result != nil {
		raw, _ := json.Marshal(result)
		var res *message.ClientResponse
		json.Unmarshal(raw, &res)
		go func() {
			m, _ := NewSendMessage("client", rID, c.id, "ack.subscription", res.MsgData)
			h.send <- m
		}()
	}
}

func (h *hub) UnsubscribeMsgNodeWs(c *client) {
	rand.Seed(time.Now().UnixNano())
	mid := rand.Int63n(999999999999)
	_, err := h.WsClient.Send(h.ctx, message.ClientSubscription{
		CID:     c.id,
		MsgId:   mid,
		MsgType: "message.unsubscribe",
		MsgData: c.msgs,
	})
	if err != nil {
		log.Println(err)
	}
}

func (h *hub) SubscribeBlockWs(c *client, rID int64) {
	rand.Seed(time.Now().UnixNano())
	mid := rand.Int63n(999999999999)
	result, err := h.WsClient.Send(h.ctx, message.ClientSubscription{
		CID:     c.id,
		MsgId:   mid,
		MsgType: "block.subscribe",
	})
	if err != nil {
		log.Println(err)
	}
	if result != nil {
		raw, _ := json.Marshal(result)
		var res *message.ClientResponse
		json.Unmarshal(raw, &res)
		go func() {
			m, _ := NewSendMessage("client", rID, c.id, "ack.subscription", res.MsgData)
			h.send <- m
		}()
	}
}

func (h *hub) UnsubscribeBlockWs(c *client) {
	rand.Seed(time.Now().UnixNano())
	mid := rand.Int63n(999999999999)
	_, err := h.WsClient.Send(h.ctx, message.ClientSubscription{
		CID:     c.id,
		MsgId:   mid,
		MsgType: "block.unsubscribe",
	})
	if err != nil {
		log.Println(err)
	}
}
