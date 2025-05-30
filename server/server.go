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
	"github.com/kovey/cli-go/util"
	"github.com/kovey/debug-go/debug"
	"github.com/kovey/discovery/etcd"
	"github.com/kovey/discovery/krpc"
	"github.com/kovey/kom"
	"github.com/kovey/kom/internal"
	"github.com/kovey/kom/service"
	"google.golang.org/grpc"
)

const (
	command_create = "create"
	arg_path       = "path"
)

type server struct {
	*app.ServBase
	e     EventInterface
	wait  sync.WaitGroup
	pprof *http.Server
}

func newServer(e EventInterface) *server {
	return &server{ServBase: &app.ServBase{}, e: e, wait: sync.WaitGroup{}}
}

func (s *server) Init(a app.AppInterface) error {
	location, err := time.LoadLocation(os.Getenv("APP_TIME_ZONE"))
	if err != nil {
		return err
	}
	time.Local = location
	if s.e != nil {
		s.e.SetName(a.Name())
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
	s.pprof = &http.Server{Addr: fmt.Sprintf("%s:%d", os.Getenv(kom.SERV_HOST), port+10001), Handler: http.DefaultServeMux}
	if err := s.pprof.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			debug.Erro("listen pprof failure, error: %s", err)
		}
	}
}

func (s *server) runTest() {
	defer s.wait.Done()
	if testOn, err := env.GetBool(kom.APP_TEST_OPEN); err != nil || !testOn {
		return
	}

	internal.Run(service.Tests())
}

func (s *server) start(a app.AppInterface) error {
	if !env.CheckDefault() {
		return fmt.Errorf(".env config not found, use create command get .env file")
	}

	if s.e != nil {
		if err := s.e.OnBefore(a); err != nil {
			return err
		}
	}

	service.Init()
	if open, err := env.GetBool(kom.APP_ETCD_OPEN); err == nil && open {
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
	}

	if s.e != nil {
		if err := s.e.OnAfter(a); err != nil {
			return err
		}
	}

	s.wait.Add(1)
	go s.runAfter()
	s.wait.Add(1)
	go s.runMonitor()
	s.wait.Add(1)
	go s.runTest()
	port, _ := env.GetInt(kom.SERV_PORT)
	debug.Info("app[%s] listen on [%s:%d]", a.Name(), os.Getenv(kom.SERV_HOST), port)
	if err := service.Listen(os.Getenv(kom.SERV_HOST), port); err != nil {
		if err == grpc.ErrServerStopped {
			return nil
		}

		debug.Erro(err.Error())
		s.Shutdown(a)
		return app.Err_Not_Restart
	}

	return nil
}

func (s *server) Run(a app.AppInterface) error {
	method, err := a.Arg(0, app.TYPE_STRING)
	if err != nil {
		method, _ = a.Get(app.Ko_Command_Start)
	}

	switch method.String() {
	case command_create:
		if s.e != nil {
			f, _ := a.Get(command_create, arg_path)
			return s.e.CreateConfig(f.String())
		}
	default:
		return s.start(a)
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

	internal.Shutdown()
	s.wait.Wait()
	return nil
}

func (s *server) Reload(a app.AppInterface) error {
	return nil
}

func (s *server) Flag(a app.AppInterface) error {
	a.FlagArg(command_create, "create config file .env")
	a.FlagLong(arg_path, util.RunDir(), app.TYPE_STRING, ".env file path created", "create")
	if s.e == nil {
		return nil
	}

	return s.e.OnFlag(a)
}

func (s *server) Usage() {
	if s.e == nil {
		s.ServBase.Usage()
		return
	}

	if !s.e.Usage() {
		s.ServBase.Usage()
	}
}
