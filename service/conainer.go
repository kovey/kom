package service

import (
	"context"

	"github.com/kovey/pool"
	"google.golang.org/grpc"
)

func container(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	nCtx := pool.NewContext(ctx)
	defer pool.PutNoCtx(nCtx)

	resp, err := handler(nCtx, req)
	return resp, err
}
