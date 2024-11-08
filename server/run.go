package server

import (
	"os"

	"github.com/kovey/cli-go/app"
	"github.com/kovey/cli-go/env"
	"github.com/kovey/debug-go/debug"
)

func Run(e EventInterface) {
	serv := newServer(e)
	cli := app.NewApp(os.Getenv(env.APP_NAME))
	cli.SetServ(serv)
	if err := cli.Run(); err != nil {
		debug.Erro(err.Error())
	}
}
