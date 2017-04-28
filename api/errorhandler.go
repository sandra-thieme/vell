package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type apiError struct {
	Error   error
	Message string
	code    int
}

type apiHandler func(w http.ResponseWriter, r *http.Request) *apiError

func (handler apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	if e := handler(w, r); e != nil {
		if e.code >= 500 {
			log.Printf("%s: %s", e.Message, e.Error)
		}
		w.WriteHeader(e.code)
		if err := json.NewEncoder(w).Encode(e); err != nil {
			panic(err)
		}
	}
}
