package main

import (
	"github.com/manasseslima/doorway/cmds"
	"github.com/manasseslima/clio"
)

func createApplication() clio.App {
	app := clio.NewApp("doorway")
	app.AddCmd("run", cmds.RunCommandHandler)
	app.AddCmd("config", cmds.ConfigCommandHandler)
	return app
}

func main() {
	app := createApplication()
	app.Run()
}