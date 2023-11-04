package server

import (
	"os"

	"github.com/kovey/discovery/etcd"
	"github.com/kovey/discovery/krpc"
	"github.com/kovey/kom/zap"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Listen krpc.Local  `yaml:"listen"` // lsiten config
	Etcd   etcd.Config `yaml:"etcd"`   // etcd config
	Zap    zap.Config  `yaml:"zap"`    // zap config
	App    App         `yaml:"app"`    // app config
}

type App struct {
	TimeZone  string `yaml:"time_zone"`
	PprofOpen string `yaml:"pprof_open"`
}

func (c *Config) Load(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal([]byte(content), c)
}
