package mem

import (
	"github.com/streamrail/concurrent-map"
)

var user_id_map map[string]int
var user_pwd_map map[string]string
var uid_token_map map[int]string

var token_map map[string]int

func init() {
	user_id_map = make(map[string]int)
	user_pwd_map = make(map[string]string)
	uid_token_map = make(map[int]string)

	token_map = make(map[string]int)
}

func Test() {
	m := cmap.New()
	m.Set("a", "bb")
	var ret interface{}
	ret, ok := m.Get("a")
	if !ok {
		panic("cmap get error")
	}
	value := ret.(string)
	if value != "bb" {
		panic("cmap error")
	}
}
