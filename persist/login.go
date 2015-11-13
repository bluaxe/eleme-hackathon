package persist

import (
	"common"
	"fmt"
)

func Login(username, passwd string) (int, bool) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	db := getDB()
	defer releaseDB(db)
	row := db.QueryRow(s.GET_USER_BY_NAME, username)
	var user common.User
	err := row.Scan(&user.Id, &user.Name, &user.Passwd)
	if err != nil {
		panic(err)
	}
	if user.Passwd == passwd {
		return user.Id, true
	} else {
		return 0, false
	}
}
