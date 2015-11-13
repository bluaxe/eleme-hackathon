package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"persist"
)

type handler func(url.Values, http.ResponseWriter)

var entrys map[string]handler = make(map[string]handler)

func init() {
	entrys["test"] = testfunc
	entrys["forbid"] = testforbid
	entrys["dbtest"] = testdb
}

func testfunc(v url.Values, w http.ResponseWriter) {
	fmt.Fprintf(w, "test func\n")
}

func testdb(v url.Values, w http.ResponseWriter) {
	fmt.Fprintf(w, persist.List10User())
}

func testforbid(v url.Values, w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
	res := &Response{
		status: http.StatusForbidden,
		Msg:    "forbidden",
	}

	str, _ := json.Marshal(res)
	fmt.Fprintf(w, string(str))
}

func dispatcher(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Form.Encode())
	en := r.Form.Get("action")
	f, have := entrys[en]
	if have && en != "" {
		f(r.Form, w)
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}
