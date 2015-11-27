package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"service"
	"sync"
	"time"
)

var order_signal *sync.Cond = sync.NewCond(&sync.Mutex{})

var local = true

var order_return_interval time.Duration = 2 * time.Second

func OrderTicker() {
	ticker := time.NewTicker(order_return_interval)
	for t := range ticker.C {
		order_signal.L.Lock()
		order_signal.Broadcast()
		order_signal.L.Unlock()
		fmt.Println(t, "Order Ticker Broad Casted.")
	}
}

func Start(addr string) {
	http.HandleFunc("/", dispatcher)
	http.HandleFunc("/login", loginDispatcher)
	http.HandleFunc("/foods", foodsDispatcher)
	http.HandleFunc("/carts", cartsDispatcher)
	http.HandleFunc("/carts/", addFood)
	http.HandleFunc("/orders", ordersDispatcher)
	http.HandleFunc("/admin/orders", adminOrdersDispatcher)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println(err)
	}
}

func writeResponse(w http.ResponseWriter, r *Response) {
	w.WriteHeader(r.status)
	ret, _ := json.Marshal(r)
	fmt.Fprintf(w, string(ret))
}

func dealRequest(w http.ResponseWriter, r *http.Request) (int, bool) {
	token := r.Header.Get("Access-Token")
	if token == "" {
		token = r.FormValue("access_token")
	}

	// fmt.Println("Debug: Recevice Reques with token: ", token)

	if token == "" {
		fmt.Println("Warning: token is empty!!!!")
		writeResponse(w, Unauthorized)
		return 0, false
	}
	var id int
	var ok bool
	if local {
		id, ok = service.CheckTokenLocal(token)
	} else {
		id, ok = service.CheckToken(token)
	}
	if !ok {
		fmt.Println("Warning: token not exist!!!!")
		writeResponse(w, Unauthorized)
		return 0, false
	}
	return id, true
}
