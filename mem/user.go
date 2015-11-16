package mem

import (
	"common"
	"fmt"
)

func SaveUser(user *common.User) {
	user_id_map[user.Name] = user.Id
	user_pwd_map[user.Name] = user.Passwd
}

func Login(username, passwd string) (int, bool) {
	defer common.RecoverAndPrint("mem Login Failed.")

	pwd, ok := user_pwd_map[username]
	if !ok {
		fmt.Println("login failed username not exist")
		return 0, false
	}
	if pwd != passwd {
		fmt.Println("login failed passwd not correct")
		return 0, false
	} else {
		id, ok := user_id_map[username]
		if ok {
			return id, ok
		}
		panic("mem get user id failed.")
	}
}
