package server

import (
	"fmt"
	"net/http"
)

func Start(addr string) {
	http.HandleFunc("/", dispatcher)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println(err)
	}
}
