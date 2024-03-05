package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"

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

func MainHandler(
	rw http.ResponseWriter,
	r *http.Request,
) {
	trans := transaction{id: uuid.New()} 
	config := cfg.Cfg 
	sd := extractserviceData(r.URL.Path, config)
	srv := config.Services[sd.Service]
	ept := srv.Endpoints[sd.Endpoint]
	url := fmt.Sprintf("%s%s%s", srv.Url, ept.Path, sd.Remaining)
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
