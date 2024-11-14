package cmds

import (
	"fmt"
	"net/http"

	cfg "msl.com/doorway/config"
	hdl "msl.com/doorway/handlers"

	"github.com/manasseslima/clio"
)

func registerHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/", hdl.MainHandler)
	mux.HandleFunc("/login", hdl.LoginHandler)
	mux.HandleFunc("/logout", hdl.LogoutHandler)
}

func RunCommandHandler(params clio.Params, values clio.Values) {
	cfg.LoadConfig(params["config"])
	server := http.NewServeMux()
	registerHandlers(server)
	addr := fmt.Sprintf(":%s", params["port"])
	http.ListenAndServe(addr, server)
}
