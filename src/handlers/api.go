package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"

	auth "msl.com/doorway/auth"
	cfg "msl.com/doorway/config"

	"github.com/google/uuid"
	"golang.org/x/exp/maps"
)

type serviceData struct {
	Service   string
	Version   string
	Endpoint  string
	Remaining string
}

type transaction struct {
	id uuid.UUID
}

type responseError struct {
	Message  string `json:"message"`
	Detail   string `json:"detail"`
	ErroCode string `json:"error-code"`
}

func extractServiceData(url string, config cfg.Config) serviceData {
	re := regexp.MustCompile(config.Pattern)
	mv := re.FindStringSubmatch(url)
	ud := serviceData{Service: "", Version: "", Endpoint: "", Remaining: ""}
	for idx, name := range re.SubexpNames() {
		switch name {
		case "service":
			ud.Service = mv[idx]
		case "version":
			ud.Version = mv[idx]
		case "endpoint":
			ud.Endpoint = mv[idx]
		case "remaining":
			ud.Remaining = mv[idx]
		}
	}
	return ud
}

func extractServiceDataFromPath(r *http.Request) serviceData {
	ud := serviceData{
		Service:   r.PathValue("service"),
		Endpoint:  r.PathValue("endpoint"),
		Version:   r.PathValue("version"),
		Remaining: r.PathValue(""),
	}
	return ud
}

func requestService(url string, method string, body []byte, headers http.Header, trans transaction) *http.Response {
	cli := http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		log.Print("Error to create new request to service")
		return nil
	}
	for _, name := range maps.Keys(headers) {
		req.Header.Set(name, headers.Get(name))
	}
	req.Header.Set("DOORWAY-TRANSACTION", trans.id.String())
	res, err := cli.Do(req)
	if err != nil {
		log.Printf("Error on requesting on service % %s: %v", method, url, err)
		return nil
	}
	return res
}

func returnError(rw http.ResponseWriter, message string, detail string, code string) {
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusForbidden)
	mess := responseError{
		Message:  message,
		Detail:   detail,
		ErroCode: code,
	}
	res, err := json.Marshal(mess)
	if err != nil {
		log.Println("Error on try response token validation error.")
	}
	rw.Write(res)
}

func MainHandler(
	rw http.ResponseWriter,
	r *http.Request,
) {
	config := cfg.Cfg
	tokenString := r.Header.Get("Authorization")
	token, err := auth.ValidateToken(tokenString)
	msgError := "Token not more valid"
	if err != nil {
		log.Println(err)
		msgError = err.Error()
	}
	if !token.Valid {
		detail := fmt.Sprintf("%s. Please renew token by login", msgError)
		returnError(rw, "Token Invalid.", detail, "0001")
		return
	} else {
		sd := extractServiceData(r.URL.Path, config)
		srv := config.Services[sd.Service]
		ept := srv.Endpoints[sd.Endpoint]
		url := fmt.Sprintf("%s/%s%s", srv.Url, ept.Path, sd.Remaining)
		trans := transaction{id: uuid.New()}
		body, _ := io.ReadAll(r.Body)
		res := requestService(url, r.Method, body, r.Header, trans)
		if res == nil {
			detail := fmt.Sprintf("Error on do request to %s", url)
			returnError(rw, "Error on do request to service", detail, "0002")
			return
		}
		defer res.Body.Close()
		rw.Header().Add("DOORWAY-TRANSACTION", trans.id.String())
		for k, v := range res.Header {
			rw.Header().Add(k, v[0])
		}
		body, err = io.ReadAll(res.Body)
		if err != nil {
			log.Print("Error to read body data")
		}
		if res.StatusCode >= 400 {
			log.Printf("[%s][Error] Request to %s -> %s: %s", trans.id.String(), r.URL.Path, url, res.Status)
		} else {
			log.Printf("[%s][Success] Request to %s -> %s: %s", trans.id.String(), r.URL.Path, url, res.Status)
		}
		rw.WriteHeader(res.StatusCode)
		rw.Write(body)
	}
}
