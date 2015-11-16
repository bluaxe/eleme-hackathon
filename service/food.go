package service

import (
	"cache"
	"common"
	"encoding/json"
	"fmt"
	"persist"
)

var food_list = func() map[int]int {
	var f = make(map[int]int)
	return f
}()

func AllFoods() string {
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
}
