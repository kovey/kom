package service

import (
	"github.com/kovey/discovery/etcd"
	"github.com/kovey/discovery/krpc"
	"github.com/kovey/discovery/register"
)

func RegisterToCenter(conf etcd.Config, ttl int64, listen *krpc.Local) error {
	register.Init(conf)
	return register.Register(listen.Instance(), ttl)
}
