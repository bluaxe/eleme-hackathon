package service

import (
	"cache"
	"common"
	"crypto/rand"
	"fmt"
)

func MakeOrder(cart_id string, uid int) string {
	cart_uid := GetCartUser(cart_id)
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
		res := FetchFood(food.Id, food.Num)
		if res < 0 {
			fmt.Printf("Warning food stock not enought fid:%d\n", food.Id)
			FetchFood(food.Id, -food.Num)
			for id, cnt := range done {
				FetchFood(id, -cnt)
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
	SaveOrder(order, uid)
	return order.Id, true
}

func GetUserOrders(uid int) *[]common.Order {
	return cache.GetUserOrders(uid)
}

func GetAllOrders() *[]common.Order {
	uids := cache.GetAllOrderUid()
	var order_list []common.Order
	for _, id := range *uids {
		var orders = *cache.GetUserOrders(id)
		for _, order := range orders {
			order_list = append(order_list, order)
		}
	}
	return &order_list
}

func SaveOrder(order *common.Order, uid int) {
	total := 0
	for _, food := range order.Foods {
		price, ok := GetFoodPrice(food.Id)
		if !ok {
			fmt.Println("Get Food Price In Mem Failed!")
		}
		total += price * food.Num
	}
	order.Total = total

	cache.SaveOrder(order, uid)
}

func NewOrderID() string {
	b := make([]byte, 1)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
