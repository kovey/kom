package service

import (
	"context"
	"fmt"

	"github.com/kovey/pool"
	"github.com/kovey/pool/object"
)

const (
	ctx_namespace = "kom.service"
	ctx_full_name = "context"
)

func init() {
	pool.Reg(pool.NewPool(ctx_namespace, ctx_full_name, func() any {
		return NewContext()
	}))
}

type Context struct {
	context.Context
	depObjs   map[string]object.ObjectInterface
	_fullName string
}

func NewContext() *Context {
	return &Context{depObjs: make(map[string]object.ObjectInterface), _fullName: fmt.Sprintf("%s.%s", ctx_namespace, ctx_full_name)}
}

func (c *Context) Init(ctx context.Context) {
	c.Context = ctx
}

func (c *Context) Reset() {
	c.Context = nil
	if len(c.depObjs) > 0 {
		c.depObjs = make(map[string]object.ObjectInterface)
	}
}

func (c *Context) FullName() string {
	return c._fullName
}

func (c *Context) Get(namespace, name string) object.ObjectInterface {
	if obj, ok := c.depObjs[fmt.Sprintf("%s.%s", namespace, name)]; ok {
		return obj
	}

	return nil
}

func (c *Context) add(obj object.ObjectInterface) {
	c.depObjs[obj.FullName()] = obj
}

func Get[T object.ObjectInterface](ctx context.Context, namespace, name string) T {
	cc := ctx.(*Context)
	val := cc.Get(namespace, name)
	return val.(T)
}
