package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"service"
)

func Start(addr string) {
	http.HandleFunc("/", dispatcher)
	http.HandleFunc("/login", loginDispatcher)
	http.HandleFunc("/foods", foodsDispatcher)
	http.HandleFunc("/carts", cartsDispatcher)
	http.HandleFunc("/carts/", addFood)
	http.HandleFunc("/orders", ordersDispatcher)
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

	fmt.Println("Recevice Reques with token: ", token)

	if token == "" {
		fmt.Println("token is empty!!!!")
		writeResponse(w, Unauthorized)
		return 0, false
	}
	id, ok := service.CheckToken(token)
	if !ok {
		fmt.Println("token not exist!!!!")
		writeResponse(w, Unauthorized)
		return 0, false
	}
	return id, true
}
