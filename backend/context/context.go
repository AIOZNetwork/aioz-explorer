package context

import (
	"context"
	cmsdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sirupsen/logrus"
)

type Context struct {
	context.Context
}

// create a new context
func NewGoBContext(height int64, logger *logrus.Entry, cmCtx cmsdk.Context) Context {
	c := Context{
		Context: context.Background(),
	}
	c = c.WithLogger(logger)
	c = c.WithBlockHeight(height)
	c = c.WithTxHash("")
	c = c.WithMsgType("")
	c = c.WithCMContext(cmCtx)
	return c
}

//----------------------------------------
// Getting a value

// context value for the provided key
func (c Context) Value(key interface{}) interface{} {
	value := c.Context.Value(key)
	return value
}

//----------------------------------------
// With* (setting a value)

// nolint
func (c Context) WithValue(key interface{}, value interface{}) Context {
	return c.withValue(key, value)
}
func (c Context) WithString(key interface{}, value string) Context {
	return c.withValue(key, value)
}
func (c Context) WithInt32(key interface{}, value int32) Context {
	return c.withValue(key, value)
}
func (c Context) WithUint32(key interface{}, value uint32) Context {
	return c.withValue(key, value)
}
func (c Context) WithUint64(key interface{}, value uint64) Context {
	return c.withValue(key, value)
}

func (c Context) withValue(key interface{}, value interface{}) Context {
	return Context{
		Context: context.WithValue(c.Context, key, value),
	}
}

//----------------------------------------
// Values that require no key.

const (
	contextKeyBlockHeight = iota
	contextKeyTxHash
	contextKeyMsgType
	contextKeyLogger
	contextKeyCMCtx
)

// --------- GET context-------------------
func (c Context) BlockHeight() int64 { return c.Value(contextKeyBlockHeight).(int64) }

func (c Context) TxHash() string { return c.Value(contextKeyTxHash).(string) }

func (c Context) MsgType() string { return c.Value(contextKeyMsgType).(string) }

func (c Context) Logger() *logrus.Entry { return c.Value(contextKeyLogger).(*logrus.Entry) }

func (c Context) CMCtx() cmsdk.Context { return c.Value(contextKeyCMCtx).(cmsdk.Context) }

// --------- SET context-------------------
func (c Context) WithBlockHeight(height int64) Context {
	c = c.WithLogger(c.Logger().WithField("blockheight", height))
	return c.withValue(contextKeyBlockHeight, height)
}

func (c Context) WithTxHash(txHash string) Context {
	return c.withValue(contextKeyTxHash, txHash)
}

func (c Context) WithMsgType(msgType string) Context {
	return c.withValue(contextKeyMsgType, msgType)
}

func (c Context) WithLogger(logger *logrus.Entry) Context {
	return c.withValue(contextKeyLogger, logger)
}

func (c Context) WithCMContext(cmctx cmsdk.Context) Context {
	return c.withValue(contextKeyCMCtx, cmctx)
}
