package server

import (
	"github.com/kovey/cli-go/app"
	"github.com/kovey/debug-go/debug"
)

func Run(e EventInterface) {
	serv := newServer(e)
	appName := "kom"
	if e != nil {
		appName = e.AppName()
	}
	cli := app.NewApp(appName)
	cli.SetServ(serv)
	if err := cli.Run(); err != nil {
		debug.Erro(err.Error())
	}
}
