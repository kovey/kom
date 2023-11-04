package service

import (
	"context"

	"github.com/kovey/pool"
	"google.golang.org/grpc"
)

func container(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	nCtx := pool.Get[*Context](ctx_namespace, ctx_full_name, ctx)
	defer put(nCtx)

	if descs, ok := svs.objs[info.FullMethod]; ok {
		for _, desc := range descs {
			if obj := pool.GetObject(desc.Namespace, desc.Name, ctx); obj != nil {
				nCtx.add(obj)
			}
		}
	}

	resp, err := handler(nCtx, req)
	return resp, err
}

func put(ctx *Context) {
	for _, obj := range ctx.depObjs {
		pool.Put(obj)
	}

	pool.Put(ctx)
}
