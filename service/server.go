package service

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
)

const (
	Net_Tcp  = "tcp"
	Net_Addr = "%s:%d"
)

var serv *grpc.Server

func Init() {
	serv = grpc.NewServer(grpc.ChainStreamInterceptor(stream_reco, stream_logger), grpc.ChainUnaryInterceptor(container, recovery, logger))
	for _, sv := range svs.svs {
		serv.RegisterService(sv.Desc(), sv)
	}

	svs.svs = nil
}

func OpenTracing(open string) {
	grpc.EnableTracing = open == "On"
}

func Register(sv ServiceInterface) {
	if serv == nil {
		svs.add(sv)
		return
	}

	serv.RegisterService(sv.Desc(), sv)
}

func Listen(host string, prot int) error {
	listener, err := net.Listen(Net_Tcp, fmt.Sprintf(Net_Addr, host, prot))
	if err != nil {
		return err
	}

	return serv.Serve(listener)
}

func Stop() {
	if serv == nil {
		return
	}

	serv.Stop()
}
