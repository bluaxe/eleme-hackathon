package mem

import (
	"common"
)

func SaveUser(user *common.User) {
	user_id_map.Set(user.Name, user.Id)
	user_pwd_map.Set(user.Name, user.Passwd)
}

func Login(username, passwd string) (int, bool) {
	defer common.RecoverAndPrint("mem Login Failed.")

	ret, ok := user_pwd_map.Get(username)
	var pwd string
	if ok {
		pwd = ret.(string)
	}

	if pwd != passwd {
		return 0, false
	} else {
		ret, ok := user_id_map.Get(username)
		if ok {
			id := ret.(int)
			return id, ok
		}
		panic("mem get user id failed.")
	}
}
