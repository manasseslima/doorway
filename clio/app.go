package clio

import (
	"os"
	"strings"
)

// 
type Command struct {
	name string
	handler func()
	commands map[string]Command
	params map[string]string
	values []string
}

func newCommand(name string, handler func()) Command {
	cmd := Command{
		name: name, 
		handler: handler,
		commands: map[string]Command{},
		params: map[string]string{},
		values: []string{},
	}
	return cmd
}

func (cmd *Command) run(args []string) {
	runme := true
	for idx, arg := range args {
		if idx == 0 && !strings.Contains(arg, "--") {
			cmd := cmd.commands[arg]
			cmd.run(args[idx + 1:])
			runme = false
			break
		}
		if strings.Index(arg, "--") == 0 {
			param := strings.Split(arg[2:], "=")
			key := param[0]
			value := param[1]
			cmd.params[key] = value
		} else {
			cmd.values = append(cmd.values, arg)
		}
	}
	if runme {
		cmd.handler()
	}
}

// Struct App manages the commands will run.
type App struct {
	name string
	commands map[string]Command
}

func NewApp(name string) App {
	app := App{name: name, commands: map[string]Command{} }
	return app
}

func (app *App) AddCmd(
	name string,
	handler func(),
) {
	cmd := newCommand(name, handler)
	app.commands[cmd.name] = cmd
}

func (app *App) Run() {
	for idx, arg := range os.Args {
		if idx == 0 {
			continue
		}
		if strings.Index(arg, "--") == 0 {
			param := strings.Split(arg[2:], "=")
			key := param[0]
			value := param[1]
			print(key, value)
		} else {
			cmd := app.commands[arg]
			cmd.run(os.Args[idx + 1:])
			break
		}
	}
}