package meta

import (
	"context"
	"encoding/json"
	"reflect"
)

type Func struct {
	Name        string
	Struct      string
	InArgs      []*Arg
	OutArgs     []*Arg
	Comment     string
	body        reflect.Value
	OutMaxIndex int
}

func (f *Func) DoHtml(q, m, version string) *DoHtml {
	return &DoHtml{M: m, Q: q, Name: f.Name, Struct: f.Struct, InArgs: f.InArgs, OutArgs: f.OutArgs, Comment: f.Comment, Version: version, OutMaxIndex: f.OutMaxIndex}
}

func (f *Func) Get(argName string) *Arg {
	for _, arg := range f.InArgs {
		if arg.ArgName == argName {
			return arg
		}
	}

	return nil
}

func (f *Func) Call(args []any) []*RunOutArg {
	inputs := make([]reflect.Value, len(args))
	for index, arg := range args {
		inputs[index] = reflect.ValueOf(arg)
	}

	tmp := f.body.Call(inputs)
	outs := make([]*RunOutArg, len(tmp))
	for i := 0; i < len(outs); i++ {
		outArg := f.OutArgs[i]
		if tmp[i].IsNil() {
			outs[i] = &RunOutArg{ArgName: outArg.ArgName, TypeName: outArg.ArgName, TypePrefix: outArg.TypePrefix, Value: "nil"}
			continue
		}

		outs[i] = &RunOutArg{ArgName: outArg.ArgName, TypeName: outArg.ArgName, TypePrefix: outArg.TypePrefix, Value: tmp[i].Interface()}
		if err, ok := outs[i].Value.(error); ok {
			outs[i].Value = err.Error()
		}
	}

	return outs
}

func (f *Func) ParseArgs(ctx context.Context, data map[string]string) ([]any, error) {
	args := make([]any, len(f.InArgs))
	for index, arg := range f.InArgs {
		if arg.TypeName == "context.Context" {
			args[index] = ctx
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
