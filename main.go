package main

import (
	"encoding/json"
	"net/http"
)

func mainHandler(
	rw http.ResponseWriter,
	r *http.Request,
) {
	data := map[string]string{
		"nome": "Manasses Lima",
		"cpf": "03518446436",
		"tel": "11984336425",
	}
	res := json.Marshal()
	rw.Write([]byte(res))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", mainHandler)
	http.ListenAndServe(":8080", mux)
}