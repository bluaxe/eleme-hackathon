package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Start(addr string) {
	http.HandleFunc("/", dispatcher)
	http.HandleFunc("/login", loginDispatcher)
	http.HandleFunc("/foods", foodsDispatcher)
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
