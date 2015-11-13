package server

import (
	"fmt"
	"net/http"
	"service"
)

func foodsDispatcher(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		writeResponse(w, BadRequest)
	}

	token := r.FormValue("access_token")
	fmt.Println(token)

	if token == "" {
		writeResponse(w, Unauthorized)
		return
	}
	_, ok := service.CheckToken(token)
	if !ok {
		writeResponse(w, Unauthorized)
		return
	}
}
