package service

import (
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
	svs []ServiceInterface
}

func newServices() *services {
	return &services{}
}

func (s *services) add(sv ServiceInterface) {
	s.svs = append(s.svs, sv)
}

var svs = newServices()
