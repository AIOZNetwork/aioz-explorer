package ws

import (
	ws "aioz.io/go-aioz/cmd/aiozmedia/io/websocket"
	gatewaytypes "aioz.io/go-aioz/cmd/aiozmedia/types"
	context2 "context"
	"encoding/json"
	"fmt"
	"github.com/tendermint/tendermint/libs/log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strings"
	"swagger-server/context"
	"swagger-server/ws/message"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	metrics "github.com/rcrowley/go-metrics"

	amino "github.com/tendermint/go-amino"
	cmn "github.com/tendermint/tendermint/libs/common"
)

const (
	defaultMaxReconnectAttempts = 25
	defaultWriteWait            = 0
	defaultReadWait             = 0
	defaultPingPeriod           = 0
)

const (
	protoHTTP  = "http"
	protoHTTPS = "https"
	protoWSS   = "wss"
	protoWS    = "ws"
	protoTCP   = "tcp"
)

// Parsed URL structure
type parsedURL struct {
	url.URL
}

// Parse URL and set defaults
func newParsedURL(remoteAddr string) (*parsedURL, error) {
	u, err := url.Parse(remoteAddr)
	if err != nil {
		return nil, err
	}

	// default to tcp if nothing specified
	if u.Scheme == "" {
		u.Scheme = protoTCP
	}

	return &parsedURL{*u}, nil
}

// Change protocol to HTTP for unknown protocols and TCP protocol - useful for RPC connections
func (u *parsedURL) SetDefaultSchemeHTTP() {
	// protocol to use for http operations, to support both http and https
	switch u.Scheme {
	case protoHTTP, protoHTTPS, protoWS, protoWSS:
		// known protocols not changed
	default:
		// default to http for unknown protocols (ex. tcp)
		u.Scheme = protoWS
	}
}

// Get full address without the protocol - useful for Dialer connections
func (u parsedURL) GetHostWithPath() string {
	// Remove protocol, userinfo and # fragment, assume opaque is empty
	return u.Host + u.EscapedPath()
}

// Get a trimmed address - useful for WS connections
func (u parsedURL) GetTrimmedHostWithPath() string {
	// replace / with . for http requests (kvstore domain)
	return strings.Replace(u.GetHostWithPath(), "/", ".", -1)
}

// Get a trimmed address with protocol - useful as address in RPC connections
func (u parsedURL) GetTrimmedURL() string {
	return u.Scheme + "://" + u.GetTrimmedHostWithPath()
}

//----------------------------------------------

func makeErrorDialer(err error) func(string, string) (net.Conn, error) {
	return func(_ string, _ string) (net.Conn, error) {
		return nil, err
	}
}

func makeHTTPDialer(remoteAddr string) func(string, string) (net.Conn, error) {
	u, err := newParsedURL(remoteAddr)
	if err != nil {
		return makeErrorDialer(err)
	}

	protocol := u.Scheme

	// accept http(s) as an alias for tcp
	switch protocol {
	case protoHTTP, protoHTTPS:
		protocol = protoTCP
	}

	return func(proto, addr string) (net.Conn, error) {
		return net.Dial(protocol, u.GetHostWithPath())
	}
}

//----------------------------------------------

type sendReq struct {
	req    gatewaytypes.IMessage
	respCh chan gatewaytypes.IResult
}

// WSClient is a WebSocket client. The methods of WSClient are safe for use by
// multiple goroutines.
type WSClient struct {
	ctx context.Context

	conn *websocket.Conn
	cdc  *amino.Codec

	idCounter uint32
	idAwaiter sync.Map

	Address  string // IP:PORT or /path/to/socket
	Endpoint string // /websocket/url/endpoint
	Queries  url.Values
	Dialer   func(string, string) (net.Conn, error)

	// Single user facing channel to read RPCResponses from, closed only when the client is being stopped.
	NotificationCh chan ws.WsResult

	// map to save wallets - device to notification
	MapWalletNotification map[string][]string

	// Callback, which will be called each time after successful reconnect.
	onReconnect  func()
	onDisconnect func()

	// internal channels
	send            chan sendReq      // user requests
	backlog         chan ws.WsMessage // stores a single user request received during a conn failure
	reconnectAfter  chan error        // reconnect requests
	readRoutineQuit chan struct{}     // a way for readRoutine to close writeRoutine

	// Maximum reconnect attempts (0 or greater; default: 25).
	maxReconnectAttempts int

	// Support both ws and wss protocols
	protocol string

	wg sync.WaitGroup

	mtx            sync.RWMutex
	sentLastPingAt time.Time
	reconnecting   bool

	// Time allowed to write a message to the server. 0 means block until operation succeeds.
	writeWait time.Duration

	// Time allowed to read the next message from the server. 0 means block until operation succeeds.
	readWait time.Duration

	// Send pings to server with this period. Must be less than readWait. If 0, no pings will be sent.
	pingPeriod time.Duration

	cmn.BaseService

	// Time between sending a ping and receiving a pong. See
	// https://godoc.org/github.com/rcrowley/go-metrics#Timer.
	PingPongLatencyTimer metrics.Timer
}

// NewWSClient returns a new client. See the commentary on the func(*WSClient)
// functions for a detailed description of how to configure ping period and
// pong wait time. The endpoint argument must begin with a `/`.
// The function panics if the provided address is invalid.
func NewWSClient(ctx context.Context, remoteAddr, endpoint string, options ...func(*WSClient)) *WSClient {
	parsedURL, err := newParsedURL(remoteAddr)
	if err != nil {
		panic(fmt.Sprintf("invalid remote %s: %s", remoteAddr, err))
	}
	// default to ws protocol, unless wss is explicitly specified
	if parsedURL.Scheme != protoWSS {
		parsedURL.Scheme = protoWS
	}

	cdc := amino.NewCodec()

	c := &WSClient{
		ctx:                  ctx,
		cdc:                  cdc,
		Address:              parsedURL.GetTrimmedHostWithPath(),
		Dialer:               makeHTTPDialer(remoteAddr),
		Endpoint:             endpoint,
		PingPongLatencyTimer: metrics.NewTimer(),

		maxReconnectAttempts:  defaultMaxReconnectAttempts,
		readWait:              defaultReadWait,
		writeWait:             defaultWriteWait,
		pingPeriod:            defaultPingPeriod,
		protocol:              parsedURL.Scheme,
		MapWalletNotification: make(map[string][]string),
	}
	c.BaseService = *cmn.NewBaseService(nil, "WSClient", c)
	for _, option := range options {
		option(c)
	}
	return c
}

// MaxReconnectAttempts sets the maximum number of reconnect attempts before returning an error.
// It should only be used in the constructor and is not Goroutine-safe.
func MaxReconnectAttempts(max int) func(*WSClient) {
	return func(c *WSClient) {
		c.maxReconnectAttempts = max
	}
}

// ReadWait sets the amount of time to wait before a websocket read times out.
// It should only be used in the constructor and is not Goroutine-safe.
func ReadWait(readWait time.Duration) func(*WSClient) {
	return func(c *WSClient) {
		c.readWait = readWait
	}
}

// WriteWait sets the amount of time to wait before a websocket write times out.
// It should only be used in the constructor and is not Goroutine-safe.
func WriteWait(writeWait time.Duration) func(*WSClient) {
	return func(c *WSClient) {
		c.writeWait = writeWait
	}
}

// PingPeriod sets the duration for sending websocket pings.
// It should only be used in the constructor - not Goroutine-safe.
func PingPeriod(pingPeriod time.Duration) func(*WSClient) {
	return func(c *WSClient) {
		c.pingPeriod = pingPeriod
	}
}

// OnReconnect sets the callback, which will be called every time after
// successful reconnect.
func OnReconnect(cb func()) func(*WSClient) {
	return func(c *WSClient) {
		c.onReconnect = cb
	}
}

// OnDisconnect sets the callback, which will be called every time after
// disconnect.
func OnDisconnect(cb func()) func(*WSClient) {
	return func(c *WSClient) {
		c.onDisconnect = cb
	}
}

func QueryString(key, value string) func(*WSClient) {
	return func(c *WSClient) {
		if c.Queries == nil {
			c.Queries = url.Values{}
		}
		c.Queries.Set(key, value)
	}
}

func Logger(logger log.Logger) func(*WSClient) {
	return func(c *WSClient) {
		c.Logger = logger.With("module", "GatewayWSClient")
	}
}

// String returns WS client full address.
func (c *WSClient) String() string {
	return fmt.Sprintf("%s (%s)", c.Address, c.Endpoint)
}

// OnStart implements cmn.Service by dialing a server and creating read and
// write routines.
func (c *WSClient) OnStart() error {
	err := c.dial()
	if err != nil {
		return err
	}

	//c.idAwaiter = make(map[uint32]chan gatewaytypes.IResult)

	c.NotificationCh = make(chan ws.WsResult, 1000)

	c.send = make(chan sendReq)
	// 1 additional error may come from the read/write
	// goroutine depending on which failed first.
	c.reconnectAfter = make(chan error, 1)
	// capacity for 1 request. a user won't be able to send more because the send
	// channel is unbuffered.
	c.backlog = make(chan ws.WsMessage, 1)

	c.startReadWriteRoutines()
	go c.reconnectRoutine()

	return nil
}

// Stop overrides cmn.Service#Stop. There is no other way to wait until Quit
// channel is closed.
func (c *WSClient) Stop() error {
	if err := c.BaseService.Stop(); err != nil {
		return err
	}
	// only close user-facing channels when we can't write to them
	c.wg.Wait()

	c.idAwaiter.Range(func(key, value interface{}) bool {
		close(value.(chan gatewaytypes.IResult))
		return true
	})
	close(c.NotificationCh)

	return nil
}

// IsReconnecting returns true if the client is reconnecting right now.
func (c *WSClient) IsReconnecting() bool {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.reconnecting
}

// IsActive returns true if the client is running and not reconnecting.
func (c *WSClient) IsActive() bool {
	return c.IsRunning() && !c.IsReconnecting()
}

// Send the given RPC request to the server. Results will be available on
// NotificationCh, errors, if any, on ErrorsCh. Will block until send succeeds or
// ctx.Done is closed.
func (c *WSClient) Send(ctx context2.Context, request gatewaytypes.IMessage) (gatewaytypes.IResult, error) {
	respCh := make(chan gatewaytypes.IResult)
	select {
	case c.send <- sendReq{req: request, respCh: respCh}:
		c.Logger.Info("sent a request", "req", request)
		resp := <-respCh
		if resp == nil {
			return nil, errors.New("request timeout")
		}
		if r, ok := resp.(gatewaytypes.ErrorResult); ok {
			return nil, errors.New(r.Error)
		}
		return resp, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (c *WSClient) ResubscribeDeviceTokens(wallets []string) {
	rand.Seed(time.Now().UnixNano())
	mid := rand.Int63n(999999999999)
	_, err := c.Send(context2.Background(), message.ClientSubscription{
		CID:     -5,
		MsgId:   mid,
		MsgType: "wallet.subscribe",
		MsgData: wallets,
	})
	if err != nil {
		fmt.Println(err)
	}
}

//// Call the given method. See Send description.
//func (c *WSClient) Call(ctx context.Context, method string, params map[string]interface{}) error {
//  request, err := types.MapToRequest(c.cdc, types.JSONRPCStringID("ws-client"), method, params)
//  if err != nil {
//    return err
//  }
//  return c.Send(ctx, request)
//}
//
//// CallWithArrayParams the given method with params in a form of array. See
//// Send description.
//func (c *WSClient) CallWithArrayParams(ctx context.Context, method string, params []interface{}) error {
//  request, err := types.ArrayToRequest(c.cdc, types.JSONRPCStringID("ws-client"), method, params)
//  if err != nil {
//    return err
//  }
//  return c.Send(ctx, request)
//}

///////////////////////////////////////////////////////////////////////////////
// Private methods

func (c *WSClient) nextID() uint32 {
	id := atomic.AddUint32(&c.idCounter, 1)
	return id
}

func (c *WSClient) newMessage(msg gatewaytypes.IMessage) ws.WsMessage {
	return ws.WsMessage{
		Id:  c.nextID(),
		Msg: msg,
	}
}

func (c *WSClient) dial() error {
	dialer := &websocket.Dialer{
		NetDial:           c.Dialer,
		Proxy:             http.ProxyFromEnvironment,
		EnableCompression: true,
	}
	queries := ""
	if len(c.Queries) > 0 {
		queries = "?" + c.Queries.Encode()
	}
	rHeader := http.Header{}
	conn, _, err := dialer.Dial(c.protocol+"://"+c.Address+c.Endpoint+queries, rHeader) // nolint:bodyclose
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

// reconnect tries to redial up to maxReconnectAttempts with exponential
// backoff.
func (c *WSClient) reconnect() error {
	attempt := 0

	c.mtx.Lock()
	c.reconnecting = true
	c.mtx.Unlock()
	defer func() {
		c.mtx.Lock()
		c.reconnecting = false
		c.mtx.Unlock()
	}()

	for {
		jitterSeconds := time.Duration(cmn.RandFloat64() * float64(time.Second)) // 1s == (1e9 ns)
		backoffDuration := jitterSeconds + ((1 << uint(attempt)) * time.Second)

		c.Logger.Info("reconnecting", "attempt", attempt+1, "backoff_duration", backoffDuration)
		time.Sleep(backoffDuration)

		err := c.dial()
		if err != nil {
			c.Logger.Error("failed to redial", "err", err)
		} else {
			c.Logger.Info("reconnected")
			if c.onReconnect != nil {
				go c.onReconnect()
			}
			return nil
		}

		attempt++

		if attempt > c.maxReconnectAttempts {
			return errors.Wrap(err, "reached maximum reconnect attempts")
		}
	}
}

func (c *WSClient) startReadWriteRoutines() {
	c.wg.Add(2)
	c.readRoutineQuit = make(chan struct{})
	go c.readRoutine()
	go c.writeRoutine()
}

func (c *WSClient) processBacklog() error {
	select {
	case msg := <-c.backlog:
		if c.writeWait > 0 {
			if err := c.conn.SetWriteDeadline(time.Now().Add(c.writeWait)); err != nil {
				c.Logger.Error("failed to set write deadline", "err", err)
			}
		}
		m, _ := json.Marshal(msg)
		if err := c.conn.WriteMessage(websocket.TextMessage, m); err != nil {
			c.Logger.Error("failed to resend request", "err", err)
			c.reconnectAfter <- err
			// requeue request
			c.backlog <- msg
			return err
		}
		c.Logger.Info("resend a request", "req", msg.Msg)
	default:
	}
	return nil
}

func (c *WSClient) reconnectRoutine() {
	for {
		select {
		case originalError := <-c.reconnectAfter:
			// wait until writeRoutine and readRoutine finish
			c.wg.Wait()
			if c.onDisconnect != nil {
				go c.onDisconnect()
			}
			if err := c.reconnect(); err != nil {
				c.Logger.Error("failed to reconnect", "err", err, "original_err", originalError)
				c.Stop()
				return
			}
			// drain reconnectAfter
		LOOP:
			for {
				select {
				case <-c.reconnectAfter:
				default:
					break LOOP
				}
			}
			err := c.processBacklog()
			if err == nil {
				c.startReadWriteRoutines()
			}

		case <-c.Quit():
			return
		}
	}
}

// The client ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *WSClient) writeRoutine() {
	var ticker *time.Ticker
	if c.pingPeriod > 0 {
		// ticker with a predefined period
		ticker = time.NewTicker(c.pingPeriod)
	} else {
		// ticker that never fires
		ticker = &time.Ticker{C: make(<-chan time.Time)}
	}

	defer func() {
		ticker.Stop()
		c.conn.Close()
		// err != nil {
		// ignore error; it will trigger in tests
		// likely because it's closing an already closed connection
		// }
		c.wg.Done()
	}()

	for {
		select {
		case req := <-c.send:
			if c.writeWait > 0 {
				if err := c.conn.SetWriteDeadline(time.Now().Add(c.writeWait)); err != nil {
					c.Logger.Error("failed to set write deadline", "err", err)
				}
			}

			msg := c.newMessage(req.req)
			respCh := make(chan gatewaytypes.IResult)
			idTimeout := time.NewTimer(10 * time.Second)
			go func() {
				select {
				case resp := <-respCh:
					req.respCh <- resp
				case <-idTimeout.C:
					c.idAwaiter.Delete(msg.Id)
					req.respCh <- nil
				}
			}()
			c.idAwaiter.Store(msg.Id, respCh)

			//fmt.Println("DEBUG MSG:", string(c.cdc.MustMarshalJSON(msg)))
			m, _ := json.Marshal(msg)
			if err := c.conn.WriteMessage(websocket.TextMessage, m); err != nil {
				c.Logger.Error("failed to send request", "err", err)
				c.reconnectAfter <- err
				// add request to the backlog, so we don't lose it
				c.backlog <- msg
				return
			}
		case <-ticker.C:
			if c.writeWait > 0 {
				if err := c.conn.SetWriteDeadline(time.Now().Add(c.writeWait)); err != nil {
					c.Logger.Error("failed to set write deadline", "err", err)
				}
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				c.Logger.Error("failed to write ping", "err", err)
				c.reconnectAfter <- err
				return
			}
			c.mtx.Lock()
			c.sentLastPingAt = time.Now()
			c.mtx.Unlock()
			c.Logger.Debug("sent ping")
		case <-c.readRoutineQuit:
			return
		case <-c.Quit():
			if err := c.conn.WriteMessage(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
			); err != nil {
				c.Logger.Error("failed to write message", "err", err)
			}
			return
		}
	}
}

// The client ensures that there is at most one reader to a connection by
// executing all reads from this goroutine.
func (c *WSClient) readRoutine() {
	defer func() {
		c.conn.Close()
		// err != nil {
		// ignore error; it will trigger in tests
		// likely because it's closing an already closed connection
		// }
		c.wg.Done()
	}()

	c.conn.SetPongHandler(func(string) error {
		// gather latency stats
		c.mtx.RLock()
		t := c.sentLastPingAt
		c.mtx.RUnlock()
		c.PingPongLatencyTimer.UpdateSince(t)

		if c.readWait > 0 {
			if err := c.conn.SetReadDeadline(time.Now().Add(c.readWait)); err != nil {
				c.Logger.Error("failed to set read deadline", "err", err)
			}
		}

		c.Logger.Debug("got pong")
		return nil
	})

	for {
		// reset deadline for every message type (control or data)
		if c.readWait > 0 {
			if err := c.conn.SetReadDeadline(time.Now().Add(c.readWait)); err != nil {
				c.Logger.Error("failed to set read deadline", "err", err)
			}
		}
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			if !websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
				return
			}

			c.Logger.Error("failed to read response", "err", err)
			close(c.readRoutineQuit)
			c.reconnectAfter <- err
			return
		}

		var response ws.WsResult
		err = json.Unmarshal(data, &response)
		if err != nil {
			c.Logger.Error("failed to parse response", "err", err, "data", string(data))
			continue
		}
		// Combine a non-blocking read on BaseService.Quit with a non-blocking write on NotificationCh to avoid blocking
		// c.wg.Wait() in c.Stop(). Note we rely on Quit being closed so that it sends unlimited Quit signals to stop
		// both readRoutine and writeRoutine
		if respCh, ok := c.idAwaiter.LoadAndDelete(response.Id); ok {
			select {
			case <-c.Quit():
			case respCh.(chan gatewaytypes.IResult) <- response.Result:
			}
		} else {
			select {
			case <-c.Quit():
			case c.NotificationCh <- response:
			}
		}

	}
}

///////////////////////////////////////////////////////////////////////////////
// Predefined methods

//// Subscribe to a query. Note the server must have a "subscribe" route
//// defined.
//func (c *WSClient) Subscribe(ctx context.Context, query string) error {
//  params := map[string]interface{}{"query": query}
//  return c.Call(ctx, "subscribe", params)
//}
//
//// Unsubscribe from a query. Note the server must have a "unsubscribe" route
//// defined.
//func (c *WSClient) Unsubscribe(ctx context.Context, query string) error {
//  params := map[string]interface{}{"query": query}
//  return c.Call(ctx, "unsubscribe", params)
//}
//
//// UnsubscribeAll from all. Note the server must have a "unsubscribe_all" route
//// defined.
//func (c *WSClient) UnsubscribeAll(ctx context.Context) error {
//  params := map[string]interface{}{}
//  return c.Call(ctx, "unsubscribe_all", params)
//}
