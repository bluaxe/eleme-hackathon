package server

import (
	"common"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"service"
	"time"
)

type order_ok struct {
	Id string `json:"id"`
}

type request_make_order struct {
	CardID string `json:"cart_id"`
}

func ordersDispatcher(w http.ResponseWriter, r *http.Request) {
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

	if r.Method == "GET" {
		orders := service.GetUserOrders(id)
		b, _ := json.Marshal(orders)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, string(b))
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	var req request_make_order
	if err := json.Unmarshal(body, &req); err != nil {
		fmt.Println("Server Make order Unmarshal Error. :", err)
		writeResponse(w, BadFormat)
		return
	}

	ret := service.MakeOrder(req.CardID, id)
	respon, have := status[ret]
	if have {
		writeResponse(w, respon)
	} else {
		fmt.Printf("Order Ok ID:%s\n", ret)
		w.WriteHeader(http.StatusOK)
		var response_ok = &order_ok{
			Id: ret,
		}
		response_ok_string, _ := json.Marshal(response_ok)
		fmt.Fprintf(w, string(response_ok_string))
	}
}
