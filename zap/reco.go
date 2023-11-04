package zap

import (
	"context"
	"fmt"

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	Err_Panic codes.Code = 500
)

func Reco(ctx context.Context, p interface{}) error {
	level := grpc_zap.DefaultCodeToLevel(Err_Panic)
	msg := fmt.Sprintf("panic: %s", p)
	ctxzap.Extract(ctx).WithOptions(zap.AddCallerSkip(6)).Check(level, msg).Write()
	return status.Errorf(Err_Panic, msg)
}

func MessageProducer(ctx context.Context, msg string, level zapcore.Level, code codes.Code, err error, duration zapcore.Field) {
	if code == Err_Panic {
		return
	}

	ctxzap.Extract(ctx).WithOptions(zap.AddCallerSkip(6)).Check(level, msg).Write(zap.Error(err))
}
