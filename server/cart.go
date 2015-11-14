package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"service"
	"strings"
)

type newCartID struct {
	Id string `json:"cart_id"`
}

type request_add_food struct {
	Food_id int `json:"food_id"`
	Count   int `json:"count"`
}

func cartsDispatcher(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		writeResponse(w, BadRequest)
		return
	}

	id, ok := dealRequest(w, r)
	if !ok {
		return
	}

	if r.Method == "POST" {
		cart := service.NewCart(id)
		var c = &newCartID{
			Id: cart,
		}

		fmt.Printf("get cart ok now return. new cart id:%s\n", cart)

		ret, _ := json.Marshal(*c)
		fmt.Fprintf(w, string(ret))
	}
}

func addFood(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Panic in Server Add Food", r)
		}
	}()
	uid, ok := dealRequest(w, r)
	if !ok {
		return
	}

	urls := strings.Split(r.URL.Path, "/")
	cart_id := urls[2]

	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	var req request_add_food
	if err := json.Unmarshal(body, &req); err != nil {
		fmt.Println("Server Cart Add Food Unmarshal Error. :", err)
		writeResponse(w, BadFormat)
		return
	}

	fmt.Printf("got request on add food, fid:%d, count : %d\n", req.Food_id, req.Count)

	res := service.AddFood(req.Food_id, req.Count, uid, cart_id)
	if res == "ok" {
		w.WriteHeader(http.StatusNoContent)
		return
	} else {
		writeResponse(w, status[res])
	}
}
