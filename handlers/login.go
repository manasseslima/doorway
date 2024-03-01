package handlers

import (
	"net/http"
)

func LoginHandler(
	rw http.ResponseWriter,
	r *http.Request,
) {
	rw.Write([]byte("login"))
}

func LogoutHandler(
	rw http.ResponseWriter,
	r *http.Request,
) {
	res := []byte("logout")
	rw.Write(res)
}