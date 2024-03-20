package auth

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	cfg "github.com/manasseslima/doorway/config"
)


type loginToken struct {
	name string
}


type ResponseLogin struct {
	Fullname string `json:"fullname"`
	Token string `json:"token"`
}

func GenerateJwt(username string) string {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp": time.Now().Add(10 * time.Minute).Unix(),
		},
	)
	tokenString, err := token.SignedString([]byte(cfg.Cfg.SecretKey))
	if err != nil {
		log.Println("Error to generate login token")
	}
	return tokenString
}