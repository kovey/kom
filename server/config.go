package server

import (
	"os"

	"github.com/kovey/discovery/etcd"
	"github.com/kovey/discovery/krpc"
	"github.com/kovey/kom/zap"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Listen   krpc.Local  `yaml:"listen"`
	Etcd     etcd.Config `yaml:"etcd"`
	Zap      zap.Config  `yaml:"zap"`
	TimeZone string      `yaml:"time_zone"`
}

func (c *Config) Load(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal([]byte(content), c)
}
