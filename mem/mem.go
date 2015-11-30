package mem

import (
	"github.com/streamrail/concurrent-map"
	"sync"
)

var user_id_map map[string]int
var user_pwd_map map[string]string
var uid_token_map map[int]*string

var token_map map[string]int

var food_stock map[int]*int32
var food_lock map[int]*sync.Mutex

var orders map[int]*string
var carts map[int]*string
var login_rets map[int]*string

func init() {
	user_id_map = make(map[string]int)
	user_pwd_map = make(map[string]string)
	uid_token_map = make(map[int]*string)

	token_map = make(map[string]int)

	food_stock = make(map[int]*int32)
	food_lock = make(map[int]*sync.Mutex)

	carts = make(map[int]*string)
	login_rets = make(map[int]*string)
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
