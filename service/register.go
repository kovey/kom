package service

import (
	"fmt"

	"github.com/kovey/discovery/etcd"
	"github.com/kovey/discovery/grpc"
	"github.com/kovey/discovery/krpc"
	"github.com/kovey/discovery/register"
)

func RegisterToCenter(conf etcd.Config, ttl int64, listen *krpc.Local) error {
	register.Init(conf)
	return register.Register(&grpc.Instance{Name: string(listen.Name), Addr: fmt.Sprintf("%s:%d", register.InnerIp(), listen.Port), Group: listen.Group, Weight: 10}, ttl)
}
