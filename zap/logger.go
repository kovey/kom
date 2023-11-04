package zap

import (
	"os"

	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func New(c Config) *zap.Logger {
	w := writer(c.Logger, c.Env)
	e := encoder(c.Env)
	core := zapcore.NewCore(e, w, level(c.Level))
	log := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	grpczap.ReplaceGrpcLoggerV2(log)
	return log
}

func level(l string) zapcore.Level {
	if lv, err := zapcore.ParseLevel(l); err == nil {
		return lv
	}

	return zap.ErrorLevel
}

func writer(w *lumberjack.Logger, env string) zapcore.WriteSyncer {
	if env == "prod" {
		return zapcore.AddSync(w)
	}

	return zapcore.AddSync(os.Stdout)
}

func encoder(env string) zapcore.Encoder {
	if env == "prod" {
		return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	}

	return zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
}
