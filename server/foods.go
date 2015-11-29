package server

import (
	// "common"
	"fmt"
	"net/http"
	"service"
	// "time"
)

func foodsDispatcher(w http.ResponseWriter, r *http.Request) {
	// now_t := time.Now()
	// defer common.LogTime(now_t, r.URL.String())

	if r.Method != "GET" {
		writeResponse(w, BadRequest)
	}
	_, ok := dealRequest(w, r)
	if !ok {
		return
	}

	ret := service.AllFoods()
	fmt.Fprintf(w, ret)
}
