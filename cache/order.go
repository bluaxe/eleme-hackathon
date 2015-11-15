package cache

import (
	"bytes"
	"common"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

func NewOrderID() (id int) {
	defer common.RecoverAndPrint("New Order ID Cache Error! ")

	c := getCon()
	defer releaseCon(c)

	key := k.ORDER_ID_KEY

	id, err := redis.Int(c.Do("incr", key))
	if err != nil {
		panic(err)
	}
	return id
}

func SaveOrder(order *common.Order, uid int) {
	defer common.RecoverAndPrint("Save Order Error! ")

	c := getCon()
	defer releaseCon(c)

	total := 0
	for _, food := range order.Foods {
		price := GetFoodPrice(food.Id)
		total += price * food.Num
	}
	order.Total = total

	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(order); err != nil {
		panic(err)
	}
	orderbytes := buf.Bytes()

	key := getUserOrderKey(uid)
	_, err := redis.Int(c.Do("hset", key, order.Id, orderbytes))
	if err != nil {
		panic(err)
	}
}

func GetAllOrderUid() *[]int {
	defer common.RecoverAndPrint("Get All Order Key Error")

	c := getCon()
	defer releaseCon(c)

	keys, err := redis.Strings(c.Do("keys", "uorder:*"))
	if err != nil {
		panic(err)
	}
	var ids []int
	for _, key := range keys {
		var id int
		fmt.Sscanf(key, "uorder:%d", &id)
		fmt.Printf("key:%s id:%d\n", key, id)
		ids = append(ids, id)
	}
	return &ids
}

func GetUserOrderLen(uid int) int {
	defer common.RecoverAndPrint("Get User Order Len")
	c := getCon()
	defer releaseCon(c)

	key := getUserOrderKey(uid)

	length, err := redis.Int(c.Do("hlen", key))
	if err != nil {
		panic(err)
	}
	return length
}

func GetUserOrders(uid int) *[]common.Order {
	defer common.RecoverAndPrint("Get User Orders Faild!")

	c := getCon()
	defer releaseCon(c)

	key := getUserOrderKey(uid)

	values, err := redis.Values(c.Do("hkeys", key))
	if err != nil {
		panic(err)
	}
	var order_ids []int
	if err := redis.ScanSlice(values, &order_ids); err != nil {
		panic(err)
	}
	// fmt.Println("Keys All Order", order_ids)

	var orders = make([]common.Order, 0)

	for _, id := range order_ids {
		orderbytes, err := redis.Bytes(c.Do("hget", key, id))
		if err != nil {
			fmt.Println("hget order failed", err, key, id)
			continue
		}

		buf := new(bytes.Buffer)
		buf.Write(orderbytes)
		dec := gob.NewDecoder(buf)

		var order common.Order
		if err = dec.Decode(&order); err != nil {
			fmt.Println("decode order failed", err)
			continue
		}
		orders = append(orders, order)
	}
	return &orders
}

func test(order common.Order) {
	defer common.RecoverAndPrint("test serial failed !")
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)

	err := enc.Encode(order)
	orderbytes := buf.Bytes()

	c := getCon()
	defer releaseCon(c)

	_, err = c.Do("set", "testkey", orderbytes)
	if err != nil {
		panic(err)
	}

	get, err := redis.Bytes(c.Do("get", "testkey"))
	if err != nil {
		panic(err)
	}
	newb := new(bytes.Buffer)
	newb.Write(get)

	dec := gob.NewDecoder(newb)
	var or common.Order
	if err = dec.Decode(&or); err != nil {
		panic(err)
	}
	s, err := json.Marshal(or)
	fmt.Println(string(s))
}
