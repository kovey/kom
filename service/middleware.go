package service

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/kovey/debug-go/debug"
	c "github.com/kovey/kom/context"
	"github.com/kovey/pool"
	"google.golang.org/grpc"
)

type logInfo struct {
	Path     string `json:"path"`
	Delay    string `json:"delay"`
	Error    string `json:"error"`
	TraceId  string `json:"trace_id"`
	SpanId   string `json:"span_id"`
	Request  string `json:"request"`
	Response string `json:"response"`
}

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

		cc, _ := ctx.(*pool.Context).Parent().(*c.Context)
		cc.Log.Erro("%s %s\r\n%s", info.FullMethod, err, stack())
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

	if !debug.FormatIsJson() {
		debug.Info("%s %s %.3fms %s", streamName, info.FullMethod, delay, errStr)
	} else {
		debug.Json(logInfo{Path: info.FullMethod, Delay: fmt.Sprintf("%.3fms", delay), Error: errStr})
	}
	return err
}

func logger(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	begin := time.Now().UnixMicro()
	resp, err = handler(ctx, req)
	delay := float64(time.Now().UnixMicro()-begin) * 0.001
	reqData, _ := json.Marshal(req)
	respDta, _ := json.Marshal(resp)
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}

	cc, _ := ctx.(*pool.Context).Parent().(*c.Context)
	if !debug.FormatIsJson() {
		cc.Log.Info("%s %.3fms %s %s, %s", info.FullMethod, delay, errStr, string(reqData), string(respDta))
	} else {
		debug.Json(logInfo{Path: info.FullMethod, Delay: fmt.Sprintf("%.3fms", delay), Error: errStr, Request: string(reqData), Response: string(respDta), TraceId: cc.TraceId(), SpanId: cc.SpanId()})
	}
	return resp, err
}
