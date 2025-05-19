package context

import (
	"context"

	"github.com/kovey/discovery/krpc"
)

type Context struct {
	context.Context
	traceId string
}

func NewContext(parent context.Context, traceId string) *Context {
	ctx := &Context{Context: context.WithValue(parent, krpc.Ko_Trace_Id, traceId), traceId: traceId}
	return ctx
}

func (c *Context) TraceId() string {
	return c.traceId
}
