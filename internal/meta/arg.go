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

type RunOutArg struct {
	ArgName    string `json:"arg_name"`
	TypeName   string `json:"type_name"`
	TypePrefix string `json:"type_prefix"`
	Value      any    `json:"value"`
}
