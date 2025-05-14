package service

import (
	c "context"

	"github.com/kovey/kom/context"
	"github.com/kovey/pool"
	"google.golang.org/grpc"
)

func container(ctx c.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	nCtx := pool.NewContext(context.NewContext(ctx))
	defer nCtx.Drop()

	return handler(nCtx, req)
}
