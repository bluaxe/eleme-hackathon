package service

import (
	"cache"
	"common"
	"fmt"
	"mem"
	"persist"
	"time"
)

func LoadAllUserToMem() {
	users := persist.GetAllUsers()
	fmt.Printf("User Num from persist: %d\n", len(*users))
	for _, user := range *users {
		mem.SaveUser(&user)
	}
}

func GenerateTokens() {
	defer common.LogTime(time.Now(), "Generate token done.")
	defer common.RecoverAndPrint("Generate Tokens Failed.")
	users := persist.GetAllUsers()
	for _, user := range *users {
		l := cache.NewLockWithExpire(user.Name, 3*1000)
		func() {
			l.GetWait()
			defer l.Release()

			token, ok := cache.UserGetToken(user.Id)
			if !ok {
				token = newToken()
				cache.UserSetToken(user.Id, token)
			}
			mem.SaveToken(token, user.Id)
			mem.UserSetToken(user.Id, token)
		}()
	}
}

func LoginLocal(username, passwd string) (int, string, bool) {
	id, ok := mem.Login(username, passwd)
	if ok {
		token := mem.UserGetToken(id)
		return id, token, true
	} else {
		return 0, "", false
	}
}

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
