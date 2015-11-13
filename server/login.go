package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"service"
)

type responseOK struct {
	Id       int    `json:"user_id"`
	Username string `json:"username"`
	Token    string `json:"access_token"`
}

type request struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var authFail = &Response{
	status: http.StatusForbidden,
	Code:   "USER_AUTH_FAIL",
	Msg:    "用户名或密码错误",
}

func loginDispatcher(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got a request !")
	if r.Method != "POST" {
		fmt.Println("Not Post!")
		writeResponse(w, BadRequest)
		return
	}

	if r.ContentLength == 0 {
		fmt.Println("Length 0")
		writeResponse(w, BadRequest)
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	fmt.Println(string(body))

	var req request
	if err := json.Unmarshal(body, &req); err != nil {
		fmt.Println(err)
		writeResponse(w, BadFormat)
		return
	}

	fmt.Printf("%s:%s\n", req.Username, req.Password)

	id, token, ok := service.Login(req.Username, req.Password)
	if !ok {
		writeResponse(w, authFail)
		return
	} else {
		var authOk = &responseOK{
			Id:       id,
			Username: req.Username,
			Token:    token,
		}
		w.WriteHeader(http.StatusOK)
		ret, _ := json.Marshal(authOk)
		fmt.Fprintf(w, string(ret))
	}
}
