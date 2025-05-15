package service

import (
	"github.com/kovey/kom/internal"
	"google.golang.org/grpc"
)

type ServiceInterface interface {
	Desc() *grpc.ServiceDesc
}

type Base struct {
	desc *grpc.ServiceDesc
}

func NewBase(desc *grpc.ServiceDesc) *Base {
	return &Base{desc: desc}
}

func (b *Base) Desc() *grpc.ServiceDesc {
	return b.desc
}

type services struct {
	svs   []ServiceInterface
	tests *internal.Serv
}

func newServices() *services {
	return &services{tests: internal.NewServ()}
}

func (s *services) add(sv ServiceInterface) {
	s.svs = append(s.svs, sv)
	s.tests.Register(sv)
}

var svs = newServices()

func Tests() *internal.Serv {
	return svs.tests
}
