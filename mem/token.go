package mem

import (
	"sync"
)

var token_map_lock sync.Mutex
var uid_token_map_lock sync.Mutex

func SaveToken(token string, id int) {
	token_map_lock.Lock()
	defer token_map_lock.Unlock()
	token_map[token] = id
}

func GetToken(token string) (int, bool) {
	id, ok := token_map[token]
	return id, ok
}

func UserGetToken(uid int) string {
	return uid_token_map[uid]
}

func UserSetToken(uid int, token string) {
	uid_token_map_lock.Lock()
	defer uid_token_map_lock.Unlock()
	uid_token_map[uid] = token
}
