package server

import (
	"os"

	"github.com/kovey/cli-go/app"
	"github.com/kovey/cli-go/env"
	"github.com/kovey/debug-go/debug"
)

func Run(e EventInterface) {
	serv := newServer(e)
	appName := "kom"
	if e != nil {
		appName = e.AppName()
	}
	cli := app.NewApp(appName)
	cli.SetDebugLevel(debug.DebugType(os.Getenv(env.DEBUG_LEVEL)))
	if line, err := env.GetInt(env.DEBUG_SHOW_FILE); err == nil {
		debug.SetFileLine(debug.FileLine(line))
	}
	cli.SetServ(serv)
	if err := cli.Run(); err != nil {
		debug.Erro(err.Error())
	}
}
