package handlers

import (
	"net/http"
)

func MainHandler(
	rw http.ResponseWriter,
	r *http.Request,
) {
	rw.Write([]byte("main"))
}
