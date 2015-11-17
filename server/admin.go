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
	defer common.LogTime(time.Now(), r.URL.String())
	defer common.RecoverAndPrint("Sever Admin Order Failed")

	if service.OverFlow() {
		return
	}

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
