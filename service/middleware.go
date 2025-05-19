package service

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/kovey/debug-go/debug"
	"github.com/kovey/discovery/krpc"
	"google.golang.org/grpc"
)

func stack() string {
	res := make([]string, 0)
	for i := 3; ; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		res = append(res, fmt.Sprintf("%s(%d)", file, line))
	}

	return strings.Join(res, "\n")
}

func stream_reco(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	defer func() {
		err := recover()
		if err == nil {
			return
		}

		streamName := "client"
		if info.IsServerStream {
			streamName = "server"
		}

		debug.Erro("%s %s %s\n%s", streamName, info.FullMethod, err, stack())
	}()

	return handler(srv, ss)
}

func recovery(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	defer func() {
		err := recover()
		if err == nil {
			return
		}

		traceId := ""
		if tmp, ok := ctx.Value(krpc.Ko_Trace_Id).(string); ok {
			traceId = tmp
		}
		debug.Erro("%s %s %s\n%s", traceId, info.FullMethod, err, stack())
	}()
	return handler(ctx, req)
}

func stream_logger(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	begin := time.Now().UnixMicro()
	err := handler(srv, ss)
	delay := float64(time.Now().UnixMicro()-begin) * 0.001
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}
	streamName := "client"
	if info.IsServerStream {
		streamName = "server"
	}

	debug.Info("%s %s %.3fms %s", streamName, info.FullMethod, delay, errStr)
	return err
}

func logger(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	begin := time.Now().UnixMicro()
	resp, err = handler(ctx, req)
	delay := float64(time.Now().UnixMicro()-begin) * 0.001
	reqData, _ := json.Marshal(req)
	respDta, _ := json.Marshal(resp)
	traceId := ""
	if tmp, ok := ctx.Value(krpc.Ko_Trace_Id).(string); ok {
		traceId = tmp
	}
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}

	debug.Info("%s %s %.3fms %s %s, %s", traceId, info.FullMethod, delay, errStr, string(reqData), string(respDta))
	return resp, err
}
