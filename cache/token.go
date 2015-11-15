package cache

import (
	"common"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

var tokenExpire = 30 //3 * 60 * 60 // 3 Hours

func SaveToken(token string, username string, id int) {
	defer common.RecoverAndPrint("Save Token Failed!")

	c := getCon()
	defer c.Close()

	key := fmt.Sprintf(k.TOKEN_KEY, token)

	_, err := c.Do("setex", key, tokenExpire, id)
	if err != nil {
		panic(err)
	}
}

func GetToken(token string) (id int, ok bool) {
	defer common.RecoverAndPrint("Get Token Failed!")

	c := getCon()
	defer c.Close()

	key := fmt.Sprintf(k.TOKEN_KEY, token)

	id, err := redis.Int(c.Do("get", key))
	if err == redis.ErrNil {
		return 0, false
	}
	if err != nil {
		panic(err)
	} else {
		return id, true
	}
}
