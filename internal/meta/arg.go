package meta

import "reflect"

type Arg struct {
	TypeName   string
	ArgName    string
	TypePrefix string
	typ        reflect.Type
	Default    string
}

func (a *Arg) Val() any {
	if a.typ.Kind() == reflect.Pointer {
		return reflect.New(a.typ.Elem()).Interface()
	}

	return reflect.New(a.typ).Interface()
}
