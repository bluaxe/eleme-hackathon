package cache

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
)

var tokenExpire = 180 //3 * 60 * 60 // 3 Hours

func SaveToken(token string, username string, id int) {
	c := getCon()
	defer c.Close()

	key := fmt.Sprintf(k.TOKEN_KEY, token)

	_, err := c.Do("setex", key, tokenExpire, id)
	if err != nil {
		panic(err)
	}
}

func GetToken(token string) (id int, ok bool) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			ok = false
		}
	}()
	c := getCon()
	defer c.Close()

	key := fmt.Sprintf(k.TOKEN_KEY, token)

	id, err := redis.Int(c.Do("get", key))
	if err != nil {
		panic(err)
	} else {
		return id, true
	}
}
