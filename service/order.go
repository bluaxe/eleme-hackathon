package service

import (
	"cache"
	"common"
	"crypto/rand"
	"fmt"
	"time"
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

	// fmt.Printf("Make order [cart:%s, uid:%d]\n", cart_id, uid)

	foods := cache.GetCartFoods(cart_id)
	var order common.Order = common.Order{
		Foods: *foods,
	}
	ret, _ := DoOrder(&order, uid)
	/*
		ret, ok := DoOrder(&order, uid)
		if ok {
			go DestroyCart(cart_id)
		}
	*/
	return ret
}

func DoOrder(order *common.Order, uid int) (string, bool) {
	var done map[int]int = make(map[int]int)
	var fetchall bool = true
	var wait chan common.CartFood = make(chan common.CartFood)
	for _, food := range order.Foods {
		go func(food common.CartFood) {
			res := FetchFood(food.Id, food.Num)
			if res < 0 {
				// fmt.Printf("Warning food stock not enought fid:%d\n", food.Id)
				go FetchFood(food.Id, -food.Num)
				wait <- common.CartFood{
					Id: -1,
				}
			} else {
				wait <- food
			}
		}(food)
	}
	for i := 0; i < len(order.Foods); i++ {
		f := <-wait
		if f.Id >= 0 {
			done[f.Id] = f.Num
		} else {
			fetchall = false
		}
	}
	if !fetchall {
		for id, cnt := range done {
			go FetchFood(id, -cnt)
		}
		return "food_not_enough", false
	}
	order.Id = NewOrderID()
	go SaveOrder(order, uid)
	return order.Id, true
}

func GetUserOrders(uid int) *[]common.Order {
	time.Sleep(1 * time.Second)
	return cache.GetUserOrders(uid)
}

func GetAllOrders() *[]common.AdminOrder {
	uids := cache.GetAllOrderUid()
	var order_list []common.AdminOrder
	for _, id := range *uids {
		var orders = *cache.GetUserOrders(id)
		for _, order := range orders {
			adminorder := common.AdminOrder{order, id}
			order_list = append(order_list, adminorder)
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
	b := make([]byte, 6)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
