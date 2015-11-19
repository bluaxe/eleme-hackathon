package service

import (
	"cache"
	"common"
	"encoding/json"
	"fmt"
	"math/rand"
	"persist"
)

var food_list = func() map[int]int {
	var f = make(map[int]int)
	return f
}()

var request_count int = 0
var request_count_thresh int = 100

var all_food_static_string string

func AllFoods() string {
	request_count += 1
	if request_count > request_count_thresh {
		return all_food_static_string
	}
	foods := cache.GetAllFoodsStock()
	var fs []common.Food = make([]common.Food, 0)
	for _, food := range *foods {
		food.Price, _ = GetFoodPrice(food.Id)
		fs = append(fs, food)
	}
	ret, _ := json.Marshal(fs)
	return string(ret)
}

func GetFoodPrice(food_id int) (int, bool) {
	price, ok := food_list[food_id]
	return price, ok
}

func InitFoodsFromPersist() {
	foods := persist.GetAllFoods()
	for _, food := range *foods {
		food_list[food.Id] = food.Price
	}

	l := cache.NewLock("master")
	l.GetWait()
	defer l.Release()

	n := cache.FoodsNum()
	if n != 0 {
		fmt.Printf("Foods Already in cache : %d\n", n)
		return
	}

	for _, food := range *foods {
		cache.FoodImport(food)
	}

	var fs []common.Food = make([]common.Food, 0)
	for _, food := range *foods {
		food.Price, _ = GetFoodPrice(food.Id)
		food.Stock -= rand.Int() % 512
		fs = append(fs, food)
	}
	ret, _ := json.Marshal(fs)
	all_food_static_string = string(ret)
	fmt.Println(all_food_static_string)
}
