package main

import (
	"github.com/manasseslima/clio"
	"msl.com/doorway/cmds"
)

func createCommands(app *clio.App) {
	app.NewCmd("run", "Run services gateway", cmds.RunCommandHandler)
	config := clio.NewCommand("config", "Generate and manage config files", cmds.ConfigCommandHandler)
	app.AddCmd(config)
}

func main() {
	app := clio.NewApp("doorway", "An APIGateway application")
	createCommands(&app)
	app.Run()
}
