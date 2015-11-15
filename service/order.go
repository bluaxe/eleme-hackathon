package service

import (
	"cache"
	"common"
	"fmt"
)

func MakeOrder(cart_id string, uid int) string {
	cart_uid := cache.GetCartUser(cart_id)
	if cart_uid == -1 {
		return "cart_not_exist"
	}
	if cart_uid != uid {
		return "cart_not_belong"
	}
	if cache.GetUserOrderLen(uid) != 0 {
		return "too_many_orders"
	}

	fmt.Printf("Make order [cart:%s, uid:%d]\n", cart_id, uid)

	foods := cache.GetCartFoods(cart_id)
	var order common.Order = common.Order{
		Foods: *foods,
	}
	ret, ok := DoOrder(&order, uid)
	if ok {
		go DestroyCart(cart_id)
	}
	return ret
}

func DoOrder(order *common.Order, uid int) (string, bool) {
	var done map[int]int = make(map[int]int)
	var fetchall bool = true
	for _, food := range order.Foods {
		res := cache.FetchFood(food.Id, food.Num)
		if res < 0 {
			go cache.FetchFood(food.Id, -food.Num)
			for id, cnt := range done {
				go cache.FetchFood(id, -cnt)
			}
			fetchall = false
			break
		}
		done[food.Id] = food.Num
	}
	if !fetchall {
		return "food_not_enough", false
	}
	order.Id = NewOrderID()
	order.Idstring = fmt.Sprintf("%d", order.Id)
	go SaveOrder(order, uid)
	return fmt.Sprintf("%d", order.Id), true
}

func GetUserOrders(uid int) *[]common.Order {
	return cache.GetUserOrders(uid)
}

func SaveOrder(order *common.Order, uid int) {
	cache.SaveOrder(order, uid)
}

func NewOrderID() int {
	return cache.NewOrderID()
}
