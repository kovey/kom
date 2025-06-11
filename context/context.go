package context

import (
	"context"

	"github.com/kovey/debug-go/debug"
	"github.com/kovey/discovery/krpc"
)

type Context struct {
	context.Context
	traceId string
	spandId string
	Log     *debug.Log
}

func NewContext(parent context.Context, traceId string) *Context {
	ctx := &Context{Context: context.WithValue(parent, krpc.Ko_Trace_Id, traceId), traceId: traceId, spandId: SpanId()}
	ctx.Log = debug.LogWith(ctx.traceId, ctx.spandId)
	return ctx
}

func (c *Context) TraceId() string {
	return c.traceId
}

func (c *Context) SpanId() string {
	return c.spandId
}
