package cache

import (
	"crypto/rand"
	"fmt"
)

type keys struct {
	TOKEN_KEY      string
	SERVER_ID      string
	LOCK_FOR_KEY   string
	RELEASE_LOCK   string
	FOOD_STOCK_KEY string
	FOOD_PRICE_KEY string
	CART_KEY_CHK   string
	CART_KEY       string
}

var k keys = InitKeys()

func InitKeys() keys {
	var k keys
	k.TOKEN_KEY = `tok:%s`
	k.SERVER_ID = `server_size`
	k.LOCK_FOR_KEY = `lock:%s`
	k.RELEASE_LOCK = `if redis.call("get",KEYS[1]) == ARGV[1] then return redis.call("del",KEYS[1]) else return 0 end`
	k.FOOD_STOCK_KEY = `fstock`
	k.FOOD_PRICE_KEY = `fprice`
	k.CART_KEY = `bas:%s`
	k.CART_KEY_CHK = `cart:%s`
	return k
}

func getLockKey(key string) string {
	return fmt.Sprintf(k.LOCK_FOR_KEY, key)
}

func getRandString() string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func getCartKey(cart_id string) string {
	return fmt.Sprintf(k.CART_KEY, cart_id)
}

func getCartKeyCHK(cart_id string) string {
	return fmt.Sprintf(k.CART_KEY_CHK, cart_id)
}
