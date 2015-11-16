package service

import (
	"cache"
	"encoding/json"
	"fmt"
	"persist"
)

var food_list = func() map[int]int {
	var f = make(map[int]int)
	return f
}()

func AllFoods() string {
	// foods := persist.GetAllFoods()
	foods := cache.GetAllFoods()
	ret, _ := json.Marshal(foods)
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
