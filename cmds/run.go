package cmds

import (
	"fmt"
	"net/http"

	"github.com/manasseslima/doorway/clio"
	hdl "github.com/manasseslima/doorway/handlers"
)

func registerHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/", hdl.MainHandler)
	mux.HandleFunc("/login", hdl.LoginHandler)
	mux.HandleFunc("/logout", hdl.LogoutHandler)
}

func RunCommandHandler(params clio.Params, values clio.Values) {
	server := http.NewServeMux()
	registerHandlers(server)
	addr := fmt.Sprintf(":%s", params["port"])
	http.ListenAndServe(addr, server)
}