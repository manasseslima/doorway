package auth

import (
	"log"
	"time"

	cfg "msl.com/doorway/config"

	"github.com/golang-jwt/jwt"
)

type loginToken struct {
	name string
}

type ResponseLogin struct {
	Fullname string `json:"fullname"`
	Token    string `json:"token"`
}

func GenerateJwt(username string) string {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(10 * time.Minute).Unix(),
		},
	)
	tokenString, err := token.SignedString([]byte(cfg.Cfg.SecretKey))
	if err != nil {
		log.Println("Error to generate login token")
	}
	return tokenString
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	if cfg.Cfg.Autorizator.Disabled {
		tokenString = GenerateJwt("anonymous")
	}
	token, err := jwt.Parse(tokenString, func(tokenString *jwt.Token) (interface{}, error) {
		return []byte(cfg.Cfg.SecretKey), nil
	})
	if err != nil {
		log.Println("Erro to verify token!")
	}
	if !token.Valid {
		log.Println("Token not valid!")
	}
	return token, err
}
