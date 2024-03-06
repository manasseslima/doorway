package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	cfg "github.com/manasseslima/doorway/config"
)

type serviceData struct {
	Service string
	Version string
	Endpoint string
	Remaining string
} 

type transaction struct {
	id uuid.UUID
}

type responseError struct {
	Message string `json:"message"`
	Detail string `json:"detail"`
	ErroCode string `json:"error-code"`
}

func extractserviceData(url string, config cfg.Config) serviceData {
	re := regexp.MustCompile(config.Pattern)
	mv := re.FindStringSubmatch(url)
	ud := serviceData{Service: "", Version: "", Endpoint: "", Remaining: ""}
	for idx, name := range re.SubexpNames() {
		switch name {
		case "service": ud.Service = mv[idx]
		case "version": ud.Version = mv[idx]
		case "endpoint": ud.Endpoint = mv[idx]
		case "remaining": ud.Remaining = mv[idx]
		}
	}
	return ud
}

func requestService(url string, method string, body io.Reader, trans transaction) *http.Response {
	cli := http.Client{}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Print("Error to create new request to service")
	}
	req.Header.Set("DOORWAY-TRANSACTION", trans.id.String())
	res, err := cli.Do(req)
	if err != nil {
		log.Print("Error on requesting on service")
	}
	return res
}

func validateToken(tokenString string) (*jwt.Token, error) {
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

func MainHandler(
	rw http.ResponseWriter,
	r *http.Request,
) {
	config := cfg.Cfg 
	tokenString := r.Header.Get("Authorization")
	token, err := validateToken(tokenString)
	msgError := "Token not more valid"
	if err != nil {
		log.Println(err)
		msgError = err.Error()
	}
	if !token.Valid {
		rw.Header().Add("Content-Type", "application/json")
		rw.WriteHeader(http.StatusForbidden)
		mess := responseError{
			Message: "Token Invalid.",
			Detail: fmt.Sprintf("%s. Please renew token by login", msgError),
			ErroCode: "0001",
		}
		res, err := json.Marshal(mess)
		if err != nil {
			log.Println("Error on try response token validation error.")
		}
		rw.Write(res)
	} else {
		sd := extractserviceData(r.URL.Path, config)
		srv := config.Services[sd.Service]
		ept := srv.Endpoints[sd.Endpoint]
		url := fmt.Sprintf("%s%s%s", srv.Url, ept.Path, sd.Remaining)
		trans := transaction{id: uuid.New()} 
		res := requestService(url, r.Method, r.Body, trans)
		rw.Header().Add("DOORWAY-TRANSACTION", trans.id.String())
		for k, v := range res.Header {
			rw.Header().Add(k, v[0])
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Print("Error to read body data")
		}
		rw.Write(body)
		log.Println(trans.id.String(), res.Status)
	}
}
