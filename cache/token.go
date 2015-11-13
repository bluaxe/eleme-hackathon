package cache

import (
	"fmt"
)

var tokenExpire = 180 //3 * 60 * 60 // 3 Hours

func SaveToken(token string, username string, id int) {
	c := getCon()
	defer c.Close()

	key := fmt.Sprintf(k.TOKEN_KEY, token)

	_, err := c.Do("setex", key, tokenExpire, string(id))
	if err != nil {
		panic(err)
	}
}

func GetToken(token string) (int, bool) {
	return 1, true
}
