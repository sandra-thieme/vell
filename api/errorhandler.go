package api

import (
	"net/http"
)

type apiError struct {
	Error   error
	Message string
	Code    int
}

type apiHandler func(w http.ResponseWriter, r *http.Request) *apiError

func (handler apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := handler(w, r); e != nil {
		http.Error(w, e.Message, e.Code)
	}
}
