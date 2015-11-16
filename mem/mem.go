package mem

import (
	"github.com/streamrail/concurrent-map"
)

var user_id_map cmap.ConcurrentMap
var user_pwd_map cmap.ConcurrentMap

func init() {
	user_id_map = cmap.New()
	user_pwd_map = cmap.New()
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
