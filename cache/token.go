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
	defer releaseCon(c)

	key := fmt.Sprintf(k.TOKEN_KEY, token)

	_, err := c.Do("setex", key, tokenExpire, id)
	if err != nil {
		panic(err)
	}
}

func UserSetToken(uid int, token string) {
	defer common.RecoverAndPrint("cache user set token failed")

	c := getCon()
	defer releaseCon(c)

	key := getUserTokenKey(uid)
	_, err := c.Do("set", key, token)
	if err != nil {
		panic(err)
	}
}

func UserGetToken(uid int) (string, bool) {
	defer common.RecoverAndPrint("cache user get token failed")

	c := getCon()
	defer releaseCon(c)

	key := getUserTokenKey(uid)

	token, err := redis.String(c.Do("get", key))
	if err != nil {
		// panic(err)
		return "", false
	}
	return token, true
}

func UserHasToken(uid int) bool {
	defer common.RecoverAndPrint("cache user has token failed")

	c := getCon()
	defer releaseCon(c)

	key := getUserTokenKey(uid)

	_, err := redis.String(c.Do("get", key))
	if err == redis.ErrNil {
		return false
	} else {
		return true
	}
}

func GetToken(token string) (id int, ok bool) {
	defer common.RecoverAndPrint("Get Token Failed!")

	c := getCon()
	defer releaseCon(c)

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
