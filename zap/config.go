package zap

import "gopkg.in/natefinch/lumberjack.v2"

type Config struct {
	Level       string             `yaml:"level"`
	Env         string             `yaml:"env"`
	Logger      *lumberjack.Logger `yaml:"logger"`
	OpenTracing string             `yaml:"open_tracing"`
}
