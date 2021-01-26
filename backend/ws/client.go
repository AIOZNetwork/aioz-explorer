package ws

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"time"
)

var (
	pingTimeout  = 6000 * time.Second // client pings server
	writeTimeout = 5 * time.Second    // server writes to client
)

type client struct {
	Namespace     string
	Id            int64
	hub           *hub
	conn          *websocket.Conn
	send          chan []byte
	done          chan bool
	wallets       map[string]bool
	msgTypes      map[string]bool
	MapUserWallet map[int64][]string
	MapUserMsg    map[int64][]string
	MapUserBlock  map[int64]bool

	unsubDone chan bool
}

func newClient(namespace string, id int64, hub *hub, conn *websocket.Conn) *client {
	return &client{
		Namespace:     namespace,
		Id:            id,
		hub:           hub,
		conn:          conn,
		send:          make(chan []byte, 1000),
		done:          make(chan bool),
		unsubDone:     make(chan bool),
		wallets:       make(map[string]bool),
		msgTypes:      make(map[string]bool),
		MapUserWallet: make(map[int64][]string),
		MapUserMsg:    make(map[int64][]string),
		MapUserBlock:  make(map[int64]bool),
	}
}

func (c *client) write(msgType int, msg []byte) error {
	c.conn.SetWriteDeadline(time.Now().Add(writeTimeout))
	return c.conn.WriteMessage(msgType, msg)
}

func (c *client) writePump() {
	defer func() {
		c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				// The hub closed channel
				logrus.Errorf("Closed by hub")
				c.write(websocket.CloseMessage, nil)
				return
			}
			if err := c.write(websocket.TextMessage, msg); err != nil {
				logrus.Errorf("Write message error %v", err)
				return
			}
		}
	}
}

func (c *client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(pingTimeout))

	// Client must ping to keep connection
	c.conn.SetPingHandler(func(string) error {
		if err := c.conn.SetReadDeadline(time.Now().Add(pingTimeout)); err != nil {
			logrus.Errorf("Set read deadline err %v", err)
		}
		if err := c.write(websocket.PongMessage, nil); err != nil {
			logrus.Errorf("Write pong error %v", err)
		}
		return nil
	})

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			logrus.Warnf("error: %v", err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logrus.Warnf("error: %v", err)
			}
			break
		}
		receiveMsg, err := NewReceiveMessage(c.Namespace, c.Id, msg)
		if err != nil {
			logrus.Warnf("Wrong format websocket message")
			continue
		}
		c.hub.Receive <- receiveMsg
	}
}
