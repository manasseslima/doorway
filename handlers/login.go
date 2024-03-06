package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	cfg "github.com/manasseslima/doorway/config"
)


type loginToken struct {
	name string
}


type responseLogin struct {
	Token string `json:"token"`
}


func generateJwt() string {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": "User Test",
			"exp": time.Now().Add(10 * time.Minute).Unix(),
		},
	)
	tokenString, err := token.SignedString([]byte(cfg.Cfg.SecretKey))
	if err != nil {
		log.Println("Error to generate login token")
	}
	return tokenString
}

func LoginHandler(
	rw http.ResponseWriter,
	r *http.Request,
) {
	data := responseLogin{
		Token: generateJwt(),
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