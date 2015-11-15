package service

import (
	"cache"
	"persist"
)

func Login(username, passwd string) (int, string, bool) {
	id, ok := persist.Login(username, passwd)
	if ok {
		var token = newToken()
		cache.SaveToken(token, username, id)
		return id, token, true
	} else {
		return 0, "", false
	}
}
