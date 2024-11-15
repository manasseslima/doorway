package main

import (
	"msl.com/doorway/cmds"

	"github.com/manasseslima/clio"
)

func createApplication() clio.App {
	app := clio.NewApp("doorway", "An apigateway application")
	app.NewCmd("run", "Run services gateway", cmds.RunCommandHandler)
	config := clio.NewCommand("config", "Generate and manage config files", cmds.ConfigCommandHandler)
	app.AddCmd(config)
	return app
}

func main() {
	app := createApplication()
	app.Run()
}
