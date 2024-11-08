package server

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"net/http"
	_ "net/http/pprof"

	"github.com/kovey/cli-go/app"
	"github.com/kovey/cli-go/env"
	"github.com/kovey/debug-go/debug"
	"github.com/kovey/discovery/etcd"
	"github.com/kovey/discovery/krpc"
	"github.com/kovey/kom"
	"github.com/kovey/kom/service"
	"github.com/kovey/kom/zap"
	"gopkg.in/natefinch/lumberjack.v2"
)

type server struct {
	*app.ServBase
	e     EventInterface
	wait  sync.WaitGroup
	pprof *http.Server
}

func newServer(e EventInterface) *server {
	return &server{e: e, wait: sync.WaitGroup{}}
}

func (s *server) Init(a app.AppInterface) error {
	location, err := time.LoadLocation(os.Getenv("APP_TIME_ZONE"))
	if err != nil {
		return err
	}
	time.Local = location

	if s.e != nil {
		if err := s.e.OnBefore(a); err != nil {
			return err
		}
	}

	maxSize, _ := env.GetInt(kom.ZAP_LOGGER_MAX_SIZE)
	maxAge, _ := env.GetInt(kom.ZAP_LOGGER_MAX_AGE)
	maxBackups, _ := env.GetInt(kom.ZAP_LOGGER_MAX_BACKUPS)
	localTime, _ := env.GetBool(kom.ZAP_LOGGER_LOCAL_TIME)
	compress, _ := env.GetBool(kom.ZAP_LOGGER_COMPRESS)
	service.Init(zap.Config{
		Level: os.Getenv(kom.ZAP_LEVEL), Env: os.Getenv(kom.ZAP_ENV),
		OpenTracing: os.Getenv(kom.ZAP_OPEN_TRACING),
		Logger: &lumberjack.Logger{
			Filename:   os.Getenv(kom.ZAP_LOGGER_FILE_NAME),
			MaxSize:    maxSize,
			MaxAge:     maxAge,
			MaxBackups: maxBackups,
			LocalTime:  localTime,
			Compress:   compress,
		},
	})

	timeout, _ := env.GetInt(kom.ETCD_TIMEOUT)
	conf := etcd.Config{
		Endpoints:   strings.Split(os.Getenv(kom.ETCD_ENDPOINTS), ","),
		DialTimeout: timeout,
		Username:    os.Getenv(kom.ETCD_USERNAME),
		Password:    os.Getenv(kom.ETCD_PASSWORD),
		Namespace:   os.Getenv(kom.ETCD_NAMESPACE),
	}
	port, _ := env.GetInt(kom.SERV_PORT)
	weight, _ := env.GetInt(kom.SERV_WEIGHT)
	local := &krpc.Local{
		Host:    os.Getenv(kom.SERV_HOST),
		Port:    port,
		Name:    krpc.ServiceName(os.Getenv(kom.SERV_NAME)),
		Group:   os.Getenv(kom.SERV_GROUP),
		Weight:  int64(weight),
		Version: os.Getenv(kom.SERV_VERSION),
	}
	ttl, _ := env.GetInt("SERV_TTL")
	if ttl <= 0 {
		ttl = 10
	}
	if err := service.RegisterToCenter(conf, int64(ttl), local); err != nil {
		return err
	}

	if s.e != nil {
		return s.e.OnAfter(a)
	}

	return nil
}

func (s *server) runAfter() {
	defer s.wait.Done()
	if s.e == nil {
		return
	}

	if err := s.e.OnRun(); err != nil {
		debug.Erro("run event.OnRun failure, error: %s", err)
	}
}

func (s *server) runMonitor() {
	defer s.wait.Done()
	if pprofOn, err := env.GetBool(kom.APP_PPROF_OPEN); err != nil || !pprofOn {
		return
	}

	port, _ := env.GetInt(kom.SERV_PORT)
	s.pprof = &http.Server{Addr: fmt.Sprintf("%s:%d", os.Getenv(kom.SERV_HOST), port+10000), Handler: http.DefaultServeMux}
	if err := s.pprof.ListenAndServe(); err != nil {
		debug.Erro("listen pprof failure, error: %s", err)
	}
}

func (s *server) Run(a app.AppInterface) error {
	s.wait.Add(1)
	go s.runAfter()
	s.wait.Add(1)
	go s.runMonitor()

	port, _ := env.GetInt(kom.SERV_PORT)
	debug.Info("app[%s] listen on [%s:%d]", a.Name(), os.Getenv(kom.SERV_HOST), port)
	if err := service.Listen(os.Getenv(kom.SERV_HOST), port); err != nil {
		return err
	}

	return nil
}

func (s *server) Shutdown(a app.AppInterface) error {
	service.Stop()
	service.Shutdown()
	if s.e != nil {
		s.e.OnShutdown()
	}
	if s.pprof != nil {
		if err := s.pprof.Shutdown(context.Background()); err != nil {
			debug.Erro("shutdown pprof failure, error: %s", err)
		}
	}

	s.wait.Wait()
	return nil
}

func (s *server) Reload(a app.AppInterface) error {
	return nil
}

func (s *server) Flag(a app.AppInterface) error {
	if s.e == nil {
		return nil
	}

	return s.e.OnFlag(a)
}
