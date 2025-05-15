package meta

import (
	"context"
	"encoding/json"
	"reflect"
)

type Func struct {
	Name    string
	Struct  string
	InArgs  []*Arg
	OutArgs []*Arg
	Comment string
	body    reflect.Value
}

func (f *Func) DoHtml(q, m, version string) *DoHtml {
	return &DoHtml{M: m, Q: q, Name: f.Name, Struct: f.Struct, InArgs: f.InArgs, OutArgs: f.OutArgs, Comment: f.Comment, Version: version}
}

func (f *Func) Get(argName string) *Arg {
	for _, arg := range f.InArgs {
		if arg.ArgName == argName {
			return arg
		}
	}

	return nil
}

func (f *Func) Call(args []any) ([]any, error) {
	inputs := make([]reflect.Value, len(args))
	for index, arg := range args {
		inputs[index] = reflect.ValueOf(arg)
	}

	tmp := f.body.Call(inputs)
	outs := make([]any, len(tmp)-1)
	for i := 0; i < len(outs); i++ {
		if tmp[i].IsNil() {
			outs[i] = nil
			continue
		}

		outs[i] = tmp[i].Interface()
	}

	last := tmp[len(outs)]
	if last.IsNil() {
		return outs, nil
	}

	return outs, last.Interface().(error)
}

func (f *Func) ParseArgs(data map[string]string) ([]any, error) {
	args := make([]any, len(f.InArgs))
	for index, arg := range f.InArgs {
		if arg.TypeName == "context.Context" {
			args[index] = context.Background()
			continue
		}
		tmpArg := arg.Val()
		if err := json.Unmarshal([]byte(data[arg.ArgName]), tmpArg); err != nil {
			return nil, err
		}

		args[index] = tmpArg
	}

	return args, nil
}
