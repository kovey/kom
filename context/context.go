package context

import "context"

type Context struct {
	context.Context
	traceId string
}

func NewContext(parent context.Context) *Context {
	ctx := &Context{Context: parent}
	if traceId, ok := parent.Value("ko_trace_id").(string); ok {
		ctx.traceId = traceId
	}

	return ctx
}

func (c *Context) TraceId() string {
	return c.traceId
}
