package cache

import (
	"common"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

var cart_expire = 60 // 1 Hour

func SaveCart(uid int, cart_id string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("save cart error! : ", r)
		}
	}()

	key := getCartKeyCHK(cart_id)
	c := getCon()
	defer releaseCon(c)

	_, err := redis.String(c.Do("setex", key, cart_expire, uid))
	if err != nil {
		panic(err)
	}
}

func DelCart(cart_id string) {
	defer common.RecoverAndPrint("Del Cart Error")
	keychk := getCartKeyCHK(cart_id)
	key := getCartKey(cart_id)

	c := getCon()
	defer releaseCon(c)
	c.Do("del", key)
	c.Do("del", keychk)
}

func GetCartUser(cid string) (id int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Get Cart User error! : ", r)
			id = -1
		}
	}()

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
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Get Cart Foods error! : ", r)
		}
	}()

	c := getCon()
	defer releaseCon(c)

	key := getCartKey(cid)

	var foods []common.CartFood
	values, err := redis.Values(c.Do("hgetall", key))
	if err != nil {
		panic(err)
	}
	if err := redis.ScanSlice(values, &foods); err != nil {
		panic(err)
	}
	return &foods
}

func CartAddFood(cid string, food_id, count int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Cart Add Food error! : ", r)
		}
	}()

	c := getCon()
	defer releaseCon(c)

	key := getCartKey(cid)

	_, err := c.Do("hincrby", key, food_id, count)
	if err != nil {
		panic(err)
	}
}
