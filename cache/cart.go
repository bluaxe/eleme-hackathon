package cache

import (
	"common"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

var cart_expire = 60 // 1 Hour

func SaveCart(uid int, cart_id string) {
	// defer common.RecoverAndPrint("Error: Cache Save Cart Error! ")

	key := getCartKeyCHK(cart_id)
	c := getCon()
	defer releaseCon(c)

	redis.String(c.Do("setex", key, cart_expire, uid))
	/*
		_, err := redis.String(c.Do("setex", key, cart_expire, uid))
		if err != nil {
			panic(err)
		}
	*/
}

func DelCart(cart_id string) {
	defer common.RecoverAndPrint("Del Cart Error")

	c := getCon()
	defer releaseCon(c)
	key := getCartKey(cart_id)
	c.Do("del", key)
	// keychk := getCartKeyCHK(cart_id)
	// c.Do("del", keychk)
}

func GetCartUser(cid string) (id int) {
	message := fmt.Sprintf("Error Cache Get Cart User Error! cid:%s\n", cid)
	defer common.RecoverPrintDo(message, func() { id = -1 })

	c := getCon()
	defer releaseCon(c)

	key := getCartKeyCHK(cid)
	id, err := redis.Int(c.Do("get", key))
	if err == redis.ErrNil {
		return -1
	}
	if err != nil {
		panic(err)
	}
	return id
}

func GetCartFoods(cid string) *[]common.CartFood {
	// defer common.Recover()

	c := getCon()
	defer releaseCon(c)

	key := getCartKey(cid)

	var foods []common.CartFood
	values, _ := redis.Values(c.Do("hgetall", key))
	/*
		if err != nil {
			panic(err)
		}
	*/
	redis.ScanSlice(values, &foods)
	/*
		if err := redis.ScanSlice(values, &foods); err != nil {
			panic(err)
		}
	*/
	return &foods
}

func CartAddFood(cid string, food_id, count int) {
	// defer common.Recover()

	c := getCon()
	defer releaseCon(c)

	key := getCartKey(cid)

	c.Do("hincrby", key, food_id, count)
	/*
		_, err := c.Do("hincrby", key, food_id, count)
		if err != nil {
			panic(err)
		}
	*/
}
