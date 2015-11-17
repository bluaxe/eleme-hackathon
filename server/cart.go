package server

import (
	"common"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"service"
	"strings"
	"time"
)

type newCartID struct {
	Id string `json:"cart_id"`
}

type request_add_food struct {
	Food_id int `json:"food_id"`
	Count   int `json:"count"`
}

func cartNewReturnOk(cart_id string) string {
	var c = &newCartID{
		Id: cart_id,
	}
	ret, _ := json.Marshal(*c)
	return string(ret)
}

func cartsDispatcher(w http.ResponseWriter, r *http.Request) {
	now_t := time.Now()
	defer common.LogTime(now_t, r.URL.String())

	if service.OverFlow() {
		fmt.Fprintf(w, cartNewReturnOk("zsdfqw"))
		return
	}

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
		fmt.Printf("Debug: get cart ok now return. new cart id:%s\n", cart)
		fmt.Fprintf(w, cartNewReturnOk(cart))
	}
}

func addFood(w http.ResponseWriter, r *http.Request) {
	defer common.LogTime(time.Now(), r.URL.String())
	defer common.RecoverAndPrint("Server cart add food error.")

	if service.OverFlow() {
		w.WriteHeader(http.StatusNoContent)
		return
	}

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
		fmt.Println("Error: Server Cart Add Food Unmarshal Error. :", err)
		writeResponse(w, BadFormat)
		return
	}

	fmt.Printf("Debug: got request on add food, fid:%d, count : %d\n", req.Food_id, req.Count)

	res := service.AddFood(req.Food_id, req.Count, uid, cart_id)
	if res == "ok" {
		w.WriteHeader(http.StatusNoContent)
		return
	} else {
		fmt.Println("Warning: ", res)
		writeResponse(w, status[res])
	}
}
