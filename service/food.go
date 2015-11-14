package service

import (
	"cache"
	"encoding/json"
	"fmt"
	"persist"
)

var food_list = func() map[int]bool {
	var f = make(map[int]bool)
	return f
}()

func AllFoods() string {
	// foods := persist.GetAllFoods()
	foods := cache.GetAllFoods()
	ret, _ := json.Marshal(foods)
	return string(ret)
}

func InitFoodsFromPersist() {
	foods := persist.GetAllFoods()
	for _, food := range *foods {
		food_list[food.Id] = true
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
