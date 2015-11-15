package server

import (
	"common"
	"encoding/json"
	"fmt"
	"net/http"
	"service"
	"time"
)

func adminOrdersDispatcher(w http.ResponseWriter, r *http.Request) {
	now_t := time.Now()
	defer common.LogTime(now_t, r.URL.String())

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Panic in Server Make Order", r)
		}
	}()

	id, ok := dealRequest(w, r)
	if !ok {
		return
	}
	if id != 1 {
		writeResponse(w, Unauthorized)
	}

	if r.Method == "GET" {
		orders := service.GetAllOrders()
		b, _ := json.Marshal(orders)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, string(b))
		return
	}
}
