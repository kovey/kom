package service

import (
	"context"

	"github.com/kovey/pool"
	"google.golang.org/grpc"
)

func container(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	nCtx := pool.NewContext(ctx)
	defer nCtx.Drop()

	return handler(nCtx, req)
}
