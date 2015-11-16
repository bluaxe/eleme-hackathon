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

func GetAllUsers() *[]common.User {
	defer common.RecoverAndPrint("persist get all user failed.")

	db := getDB()
	defer releaseDB(db)

	rows, err := db.Query(s.SELECT_ALL_USERS)
	if err != nil {
		panic(err)
	}
	var users []common.User = make([]common.User, 0)
	for rows.Next() {
		var user common.User
		if err := rows.Scan(&user.Id, &user.Name, &user.Passwd); err != nil {
			fmt.Println("Scan user failed")
		}
		users = append(users, user)
	}

	// var users []common.User = make([]common.User, 0)
	return &users
}
