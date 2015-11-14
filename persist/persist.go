package persist

import (
	"common"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "mysql"
	"os"
)

var default_dsn = getDefaultDsn()
var s sqls = initSQL()

var default_db = func() *sql.DB {
	db, err := sql.Open("mysql", default_dsn)
	if err != nil {
		panic(err)
	}
	return db
}()

func getDefaultDsn() string {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	db := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, db)
	fmt.Printf("Default DSN:%s\n", dsn)
	return dsn
}

func getDB() *sql.DB {
	return default_db
}

func releaseDB(db *sql.DB) {
	// db.Close()
}

func List10User() string {
	func() {
		if r := recover(); r != nil {
			fmt.Println("panic in persist:", r)
		}
	}()
	db := getDB()
	defer releaseDB(db)
	err := db.Ping()
	if err != nil {
		panic("db not alive")
	}

	rows, err := db.Query(s.GET_USER_TEST)
	defer rows.Close()
	if err != nil {
		panic("user query failed!")
	}
	fmt.Printf(s.GET_USER_TEST)
	var users []common.User
	for rows.Next() {
		var user common.User
		err := rows.Scan(&user.Name, &user.Passwd)
		if err != nil {
			panic("scan error")
		}
		users = append(users, user)
	}
	str, _ := json.Marshal(users)
	return string(str)
}
