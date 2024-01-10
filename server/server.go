package server

import (
	"context"
	"fmt"
	"sync"
	"time"

	"net/http"
	_ "net/http/pprof"

	"github.com/kovey/cli-go/app"
	"github.com/kovey/debug-go/debug"
	"github.com/kovey/kom/service"
)

type server struct {
	*app.ServBase
	conf  *Config
	e     EventInterface
	wait  sync.WaitGroup
	pprof *http.Server
}

func newServer(e EventInterface) *server {
	return &server{e: e, wait: sync.WaitGroup{}}
}

func (s *server) loadConf(a app.AppInterface) error {
	path, err := a.Get("c")
	if err != nil {
		return err
	}
	tmp := path.String()
	if tmp == "" {
		return fmt.Errorf("path is emty")
	}

	conf := &Config{}
	if err := conf.Load(tmp); err != nil {
		return err
	}

	s.conf = conf
	return nil
}

func (s *server) Init(a app.AppInterface) error {
	if err := s.loadConf(a); err != nil {
		return err
	}
	location, err := time.LoadLocation(s.conf.App.TimeZone)
	if err != nil {
		return err
	}
	time.Local = location

	if s.e != nil {
		if err := s.e.OnBefore(a); err != nil {
			return err
		}
	}

	service.Init(s.conf.Zap)
	if err := service.RegisterToCenter(s.conf.Etcd, 10, &s.conf.Listen); err != nil {
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
	if s.conf.App.PprofOpen != "On" {
		return
	}

	s.pprof = &http.Server{Addr: fmt.Sprintf("%s:%d", s.conf.Listen.Host, s.conf.Listen.Port+10000), Handler: http.DefaultServeMux}
	if err := s.pprof.ListenAndServe(); err != nil {
		debug.Erro("listen pprof failure, error: %s", err)
	}
}

func (s *server) Run(a app.AppInterface) error {
	s.wait.Add(1)
	go s.runAfter()
	s.wait.Add(1)
	go s.runMonitor()

	debug.Info("app[%s] listen on [%s]", a.Name(), s.conf.Listen.Addr())
	if err := service.Listen(s.conf.Listen.Host, s.conf.Listen.Port); err != nil {
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
