package service

import (
	c "context"

	"github.com/kovey/discovery/krpc"
	"github.com/kovey/kom/context"
	"github.com/kovey/pool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func container(ctx c.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	traceId := ""
	md, ok := metadata.FromIncomingContext(ctx)
	if ok && len(md[krpc.Ko_Trace_Id]) > 0 {
		traceId = md[krpc.Ko_Trace_Id][0]
	}

	nCtx := pool.NewContext(context.NewContext(ctx, traceId))
	defer nCtx.Drop()

	return handler(nCtx, req)
}
