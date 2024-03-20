package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"github.com/manasseslima/doorway/auth"
)

func LoginHandler(
	rw http.ResponseWriter,
	r *http.Request,
) {
	username := "usertest"
	fullName := "User Test"
	data := auth.ResponseLogin{
		Fullname: fullName,
		Token: auth.GenerateJwt(username),
	}
	res, err := json.Marshal(data)
	if err != nil {
		log.Println("Error on marshal login result.")
	}
	rw.Header().Add("Content-Type", "application/json")
	rw.Header().Add("Content-Length", strconv.Itoa(len(res)))
	rw.Write(res)
}

func LogoutHandler(
	rw http.ResponseWriter,
	r *http.Request,
) {
	res := []byte("logout")
	rw.Write(res)
}