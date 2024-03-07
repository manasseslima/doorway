package main

import (
	"github.com/manasseslima/doorway/cmds"
	"github.com/manasseslima/clio"
)

func createApplication() clio.App {
	app := clio.NewApp("doorway", "An apiwateway application")
	app.AddCmd("run", "Run services gateway", cmds.RunCommandHandler)
	app.AddCmd("config", "Generate and manage config files", cmds.ConfigCommandHandler)
	return app
}

func main() {
	app := createApplication()
	app.Run()
}