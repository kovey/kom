package server

import "github.com/kovey/cli-go/app"

type EventInterface interface {
	OnFlag(app.AppInterface) error
	OnBefore(app.AppInterface) error
	OnRun() error
	OnAfter(app.AppInterface) error
	OnShutdown()
	CreateConfig(path string) error
	Usage() bool
}
