package service

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	kz "github.com/kovey/kom/zap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func stream(logger *zap.Logger) grpc.ServerOption {
	return grpc.StreamInterceptor(
		grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_opentracing.StreamServerInterceptor(),
			grpc_zap.StreamServerInterceptor(logger, grpc_zap.WithMessageProducer(kz.MessageProducer)),
			grpc_recovery.StreamServerInterceptor(grpc_recovery.WithRecoveryHandlerContext(kz.Reco)),
		),
	)
}

func unary(logger *zap.Logger) grpc.ServerOption {
	return grpc.UnaryInterceptor(
		grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_opentracing.UnaryServerInterceptor(),
			grpc_zap.UnaryServerInterceptor(logger, grpc_zap.WithMessageProducer(kz.MessageProducer)),
			grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandlerContext(kz.Reco)),
			container,
		),
	)
}
