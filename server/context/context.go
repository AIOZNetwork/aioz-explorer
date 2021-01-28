package context

import (
	"context"
	"github.com/sirupsen/logrus"
)

type Context struct {
	context.Context
}

// create a new context
func NewServerContext(logger *logrus.Entry) Context {
	c := Context{
		Context: context.Background(),
	}
	c = c.WithLogger(logger)
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
	contextKeyLogger = iota
)

// --------- GET context-------------------

func (c Context) Logger() *logrus.Entry { return c.Value(contextKeyLogger).(*logrus.Entry) }

// --------- SET context-------------------

func (c Context) WithLogger(logger *logrus.Entry) Context {
	return c.withValue(contextKeyLogger, logger)
}
