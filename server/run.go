package server

import (
	"github.com/kovey/cli-go/app"
	"github.com/kovey/debug-go/debug"
)

func Run(e EventInterface, name string) {
	serv := newServer(e)
	cli := app.NewApp(name)
	cli.SetServ(serv)
	cli.SetDebugLevel(debug.Debug_Info)
	cli.Flag("c", "", app.TYPE_STRING, "app config file path")
	if err := cli.Run(); err != nil {
		debug.Erro(err.Error())
	}
}
