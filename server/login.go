package server

import (
	"common"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mem"
	"net/http"
	"persist"
	"service"
	"time"
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

func InitLoginRetStrings() {
	users := persist.GetAllUsers()
	for _, user := range *users {
		token := mem.UserGetToken(user.Id)
		var authOk = &responseOK{
			Id:       user.Id,
			Username: user.Name,
			Token:    *token,
		}
		ret, _ := json.Marshal(authOk)
		mem.SetLoginRetString(user.Id, string(ret))
	}
}

func loginDispatcher(w http.ResponseWriter, r *http.Request) {
	now_t := time.Now()
	defer common.LogTime(now_t, r.URL.String())

	if r.Method != "POST" {
		fmt.Println("Not Post!")
		writeResponse(w, BadRequest)
		return
	}

	if r.ContentLength == 0 {
		// fmt.Println("Warning: Length 0")
		writeResponse(w, BadRequest)
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	// fmt.Println(string(body))

	var req request
	if err := json.Unmarshal(body, &req); err != nil {
		// fmt.Println(err)
		writeResponse(w, BadFormat)
		return
	}

	// fmt.Printf("Debug: Got Login request %s:%s\n", req.Username, req.Password)

	var id int
	// var token string
	var ok bool
	if local {
		id, _, ok = service.LoginLocal(req.Username, req.Password)
	} else {
		id, _, ok = service.Login(req.Username, req.Password)
	}
	if !ok {
		// fmt.Println("Warning: authFail.")
		writeResponse(w, authFail)
		return
	} else {
		/*
			var authOk = &responseOK{
				Id:       id,
				Username: req.Username,
				Token:    token,
			}
			ret, _ := json.Marshal(authOk)
		*/
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, *mem.GetLoginRetString(id))
	}
}
