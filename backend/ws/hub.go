package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
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
	mu         sync.Mutex
	Clients    map[*client]bool
	register   chan *client
	unregister chan *client
	Send       chan *sendMessage
	Receive    chan *receiveMessage
	ack        *ack

	Subscribers map[string]bool
}

func newHub() *hub {
	return &hub{
		Clients:    make(map[*client]bool),
		register:   make(chan *client),
		unregister: make(chan *client),
		Send:       make(chan *sendMessage, 1000),
		Receive:    make(chan *receiveMessage, 1000),
		ack:        newAck(),

		Subscribers: make(map[string]bool),
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
	for c, _ := range h.Clients {
		if c.Id == id && c.Namespace == namespace {
			return c
		}
	}
	return nil
}

func (h *hub) findClientsOfNamespace(namespace string) []*client {
	result := make([]*client, 0)
	for c, _ := range h.Clients {
		if namespace == "all" || c.Namespace == namespace {
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
			if old := h.findClient(c.Namespace, c.Id); old != nil {
				close(old.send)
				delete(h.Clients, old)
			} else {
				fmt.Printf("Connect %d\n", c.Id)
			}
			h.Clients[c] = true
			h.mu.Unlock()
		case c := <-h.unregister:
			if _, ok := h.Clients[c]; ok {
				h.mu.Lock()
				close(c.send)
				delete(h.Clients, c)
				h.mu.Unlock()
			}
			// todo: should remove from h.Subscribers,
			// must check if there is no client subscribed to the same topic
			// theoretically, there is only 1 client connected to ws, so feel free to delete all subscribers
			// will do nothing to subscription mechanism

		case msg := <-h.Send:
			for c := range h.Clients {
				if msg.Namespace == c.Namespace && msg.ClientId == c.Id {
					body, _ := json.Marshal(msg)
					select {
					case c.send <- body:
					default:
						h.mu.Lock()
						close(c.send)
						delete(h.Clients, c)
						h.mu.Unlock()
					}
				}
			}
		case msg := <-h.Receive:
			switch msg.MsgType {
			case "wallet.subscribe":
				client := h.findClient(msg.Namespace, msg.ClientID)
				jsonobj, _ := json.Marshal(msg.MsgData)
				data := make([]string, 0)
				_ = json.Unmarshal(jsonobj, &data)
				for _, d := range data {
					if !client.wallets[d] {
						client.wallets[d] = true
					}
					if !h.Subscribers[d] {
						h.Subscribers[d] = true
					}
				}
				client.MapUserWallet[msg.CID] = data
				go func() {
					msg, _ := NewSendMessage(client.Namespace, msg.Id, client.Id, -1, "response", "wallet subscribed successful")
					h.Send <- msg
				}()
			case "wallet.unsubscribe":
				client := h.findClient(msg.Namespace, msg.ClientID)
				delete(client.MapUserWallet, msg.CID)
				go func() {
					msg, _ := NewSendMessage(client.Namespace, msg.Id, client.Id, -1, "response", "wallet unsubscribed successful")
					h.Send <- msg
				}()
			case "message.subscribe":
				client := h.findClient(msg.Namespace, msg.ClientID)
				jsonobj, _ := json.Marshal(msg.MsgData)
				data := make([]string, 0)
				_ = json.Unmarshal(jsonobj, &data)
				for _, d := range data {
					if !client.msgTypes[d] {
						client.msgTypes[d] = true
					}
					if !h.Subscribers[d] {
						h.Subscribers[d] = true
					}
				}
				client.MapUserMsg[msg.CID] = data
				go func() {
					msg, _ := NewSendMessage(client.Namespace, msg.Id, client.Id, -1, "response", "message subscribed successful")
					h.Send <- msg
				}()
			case "message.unsubscribe":
				client := h.findClient(msg.Namespace, msg.ClientID)
				delete(client.MapUserMsg, msg.CID)
				go func() {
					msg, _ := NewSendMessage(client.Namespace, msg.Id, client.Id, -1, "response", "message unsubscribed successful")
					h.Send <- msg
				}()
			case "block.subscribe":
				client := h.findClient(msg.Namespace, msg.ClientID)
				client.MapUserBlock[msg.CID] = true
				go func() {
					msg, _ := NewSendMessage(client.Namespace, msg.Id, client.Id, -1, "response", "block subscribed successful")
					h.Send <- msg
				}()
			case "block.unsubscribe":
				client := h.findClient(msg.Namespace, msg.ClientID)
				delete(client.MapUserBlock, msg.CID)
				go func() {
					msg, _ := NewSendMessage(client.Namespace, msg.Id, client.Id, -1, "response", "block unsubscribed successful")
					h.Send <- msg
				}()
			default:
				log.Println(msg)
			}
		}
	}
}

func (h *hub) FindClientOfWallet(wallet string) []*client {
	result := make([]*client, 0)
	for c, _ := range h.Clients {
		if c.wallets[wallet] {
			result = append(result, c)
		}
	}
	return result
}

func (h *hub) FindClientOfMsgType(msgType string) []*client {
	result := make([]*client, 0)
	for c, _ := range h.Clients {
		if c.msgTypes[msgType] {
			result = append(result, c)
		}
	}
	return result
}
