package internal

import (
	"runtime"

	"github.com/kovey/kom/internal/meta"
)

type Serv struct {
	Structs []*meta.Struct
	Version string
}

func NewServ() *Serv {
	s := &Serv{Version: runtime.Version()}
	return s
}

func (s *Serv) Register(desc meta.ServiceInterface) {
	s.Structs = append(s.Structs, meta.NewStruct(desc))
}

func (s *Serv) Get(servName string) *meta.Struct {
	for _, str := range s.Structs {
		if str.Name == servName {
			return str
		}
	}

	return nil
}
