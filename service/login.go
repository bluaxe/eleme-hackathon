package service

import (
	"cache"
	"mem"
	"persist"
)

func LoadAllUserToMem() {
	users := persist.GetAllUsers()
	for _, user := range *users {
		mem.SaveUser(&user)
	}
}

func Login(username, passwd string) (int, string, bool) {
	// id, ok := persist.Login(username, passwd)
	id, ok := mem.Login(username, passwd)
	if ok {
		var token = newToken()
		cache.SaveToken(token, username, id)
		go mem.SaveToken(token, id)
		return id, token, true
	} else {
		return 0, "", false
	}
}
