package server

import (
	"common"
	"fmt"
	"net/http"
	"service"
	"time"
)

func foodsDispatcher(w http.ResponseWriter, r *http.Request) {
	now_t := time.Now()
	defer common.LogTime(now_t, r.URL.String())

	if r.Method != "GET" {
		writeResponse(w, BadRequest)
	}

	token := r.Header.Get("Access-Token")
	if token == "" {
		token = r.FormValue("access_token")
	}

	// fmt.Println(token)

	if token == "" {
		writeResponse(w, Unauthorized)
		return
	}
	_, ok := service.CheckToken(token)
	if !ok {
		writeResponse(w, Unauthorized)
		return
	}

	ret := service.AllFoods()
	fmt.Fprintf(w, ret)
}
