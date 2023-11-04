package service

import (
	"fmt"

	"google.golang.org/grpc"
)

type ServiceInterface interface {
	Desc() *grpc.ServiceDesc
	Deps() map[string][]*ObjDesc
}

type ObjDesc struct {
	Namespace string
	Name      string
	Method    string
}

func NewObjDesc(namespace, name, method string) *ObjDesc {
	return &ObjDesc{Namespace: namespace, Name: name, Method: method}
}

type method struct {
	fullName string
	name     string
}

type Base struct {
	desc *grpc.ServiceDesc
	deps map[string][]*ObjDesc
}

func getFullMethod(desc *grpc.ServiceDesc) []*method {
	res := make([]*method, len(desc.Methods))
	for index, m := range desc.Methods {
		res[index] = &method{fullName: fmt.Sprintf("/%s/%s", desc.ServiceName, m.MethodName), name: m.MethodName}
	}

	return res
}

func NewBase(desc *grpc.ServiceDesc, objs ...*ObjDesc) *Base {
	b := &Base{desc: desc, deps: make(map[string][]*ObjDesc)}
	for _, method := range getFullMethod(desc) {
		for _, obj := range objs {
			if obj.Method == "" {
				b.deps[method.fullName] = append(b.deps[method.fullName], obj)
				continue
			}

			if obj.Method == method.name {
				b.deps[method.fullName] = append(b.deps[method.fullName], obj)
			}
		}
	}

	return b
}

func (b *Base) Desc() *grpc.ServiceDesc {
	return b.desc
}

func (b *Base) Deps() map[string][]*ObjDesc {
	return b.deps
}

type services struct {
	svs  []ServiceInterface
	objs map[string][]*ObjDesc
}

func newServices() *services {
	return &services{objs: make(map[string][]*ObjDesc)}
}

func (s *services) add(sv ServiceInterface) {
	s.svs = append(s.svs, sv)
	for method, objs := range sv.Deps() {
		s.objs[method] = objs
	}
}

var svs = newServices()
