package main

import (
	"net/http"
	"github.com/manasseslima/doorway/clio"
)

func loadConfigs() {

}

func registerHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/", MainHandler)
	mux.HandleFunc("/login", LoginHandler)
	mux.HandleFunc("/logout", LogoutHandler)
}

func runCommandHandler() {
	print("command run")
}

func configCommandHandler() {
	print("command config")
}

func createApplication() clio.App {
	app := clio.NewApp("doorway")
	app.AddCmd("run", runCommandHandler)
	app.AddCmd("config", configCommandHandler)
	return app
}

func main() {
	app := createApplication()
	app.Run()
}