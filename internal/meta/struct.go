package meta

import (
	"fmt"
	"reflect"
	"strings"

	"google.golang.org/grpc"
)

type ServiceInterface interface {
	Desc() *grpc.ServiceDesc
}

type Struct struct {
	Name     string
	Short    string
	Comment  string
	Funcs    []*Func
	srvValue reflect.Value
}

func NewStruct(srv ServiceInterface) *Struct {
	return _parse(srv)
}

func (s *Struct) Get(funcName string) *Func {
	for _, fn := range s.Funcs {
		if fn.Name == funcName {
			return fn
		}
	}

	return nil
}

func (s *Struct) Call(funcName string, args []any) ([]any, error) {
	for _, fn := range s.Funcs {
		if fn.Name == funcName {
			return fn.Call(args)
		}
	}

	return nil, fmt.Errorf("method[%s] not found", funcName)
}

func _parse(desc ServiceInterface) *Struct {
	srv := desc.Desc()
	s := &Struct{Name: srv.ServiceName, Funcs: make([]*Func, len(srv.Methods)), srvValue: reflect.ValueOf(desc)}
	s.Short = strings.ToLower(string([]byte{s.Name[0]}))
	for i, m := range srv.Methods {
		s.Funcs[i] = &Func{Name: m.MethodName, Struct: srv.ServiceName}
	}

	vType := reflect.TypeOf(desc)
	for i := 0; i < vType.NumMethod(); i++ {
		fn := vType.Method(i)
		fnInfo := s.Get(fn.Name)
		if fnInfo == nil {
			continue
		}

		fnInfo.body = s.srvValue.Method(i)
		for j := 1; j < fn.Type.NumIn(); j++ {
			iType := fn.Type.In(j)
			fnInfo.InArgs = append(fnInfo.InArgs, _parseArg(iType, j, true))
		}
		for j := 0; j < fn.Type.NumOut(); j++ {
			iType := fn.Type.Out(j)
			fnInfo.OutArgs = append(fnInfo.OutArgs, _parseArg(iType, j, false))
		}
	}

	return s
}

func _parseArg(typ reflect.Type, index int, in bool) *Arg {
	a := &Arg{TypeName: typ.String(), ArgName: fmt.Sprintf("arg%d", index), typ: typ, Default: "{}"}
	if !in {
		a.ArgName = ""
	}
	switch typ.Kind() {
	case reflect.Array, reflect.Slice:
		a.TypePrefix = "[]"
		a.TypeName = typ.Elem().String()
		a.Default = "[]"
	case reflect.Map:
		a.TypePrefix = fmt.Sprintf("map[%s]", typ.Key().String())
		a.TypeName = typ.Elem().String()
	}
	return a
}

func FirstUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func FirstLower(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}
