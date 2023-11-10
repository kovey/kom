package service

import (
	"fmt"

	"google.golang.org/grpc"
)

type ServiceInterface interface {
	Desc() *grpc.ServiceDesc
}

type method struct {
	fullName string
	name     string
}

type Base struct {
	desc *grpc.ServiceDesc
}

func getFullMethod(desc *grpc.ServiceDesc) []*method {
	res := make([]*method, len(desc.Methods))
	for index, m := range desc.Methods {
		res[index] = &method{fullName: fmt.Sprintf("/%s/%s", desc.ServiceName, m.MethodName), name: m.MethodName}
	}

	return res
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
