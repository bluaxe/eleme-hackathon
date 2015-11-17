package service

import (
	"cache"
	"common"
	"fmt"
	"mem"
	"persist"
	"sync"
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
	var wg sync.WaitGroup

	var gen = func(us []common.User, key string) {
		defer common.LogTime(time.Now(), key)
		defer common.RecoverAndPrint(key)
		defer wg.Done()
		l := cache.NewLockWithExpire(key, 300*1000)
		l.GetWait()
		defer l.Release()
		_, has := cache.UserGetToken(us[0].Id)
		if !has {
			fmt.Println("I am generating.", key)
		} else {
			l.Release()
		}
		for _, user := range us {
			var token string
			var ok bool
			if has {
				token, ok = cache.UserGetToken(user.Id)
				if !ok {
					panic("not found token")
				}
			} else {
				token = newToken()
				cache.UserSetToken(user.Id, token)
			}
			mem.SaveToken(token, user.Id)
			mem.UserSetToken(user.Id, token)
			fmt.Printf("Token Gene: %d :%s:%s\n", user.Id, user.Name, token)
		}
	}

	length := len(*users)
	step := 10000
	for start := 0; start < length; start += step {
		key := fmt.Sprintf("tokengen%d", start)
		cut := (*users)[start : start+step]
		wg.Add(1)
		go gen(cut, key)
	}
	wg.Wait()

	/*

		for _, user := range *users {
			l := cache.NewLockWithExpire(user.Name, 5000)
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
				// fmt.Printf("Token Gene:%d:%s:%s\n", user.Id, user.Name, token)
			}()
		}
	*/

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
