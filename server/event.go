package server

import (
	"os"

	"github.com/kovey/cli-go/app"
	"github.com/kovey/cli-go/env"
)

type EventInterface interface {
	OnFlag(app.AppInterface) error
	OnBefore(app.AppInterface) error
	OnRun() error
	OnAfter(app.AppInterface) error
	OnShutdown()
	CreateConfig(path string) error
	Usage() bool
	SetName(name string)
	AppName() string
	Version() string
	Author() string
}

type EventBase struct {
	name string
}

func (s *EventBase) SetName(name string) {
	s.name = name
}

func (s *EventBase) OnBefore(app.AppInterface) error {
	return nil
}

func (s *EventBase) OnAfter(app.AppInterface) error {
	return nil
}

func (s *EventBase) OnRun() error {
	return nil
}

func (s *EventBase) OnShutdown() {
}

func (s *EventBase) CreateConfig(app.AppInterface) error {
	return nil
}

func (s *EventBase) Usage() bool {
	return false
}

func (s *EventBase) AppName() string {
	return os.Getenv(env.APP_NAME)
}

func (s *EventBase) Version() string {
	return ""
}

func (s *EventBase) Author() string {
	return ""
}
