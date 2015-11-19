package service

import (
	"cache"
	"common"
	"encoding/json"
	"fmt"
	"math/rand"
	"mem"
	"persist"
	"sync"
)

var food_list = func() map[int]int {
	var f = make(map[int]int)
	return f
}()

var request_count int = 0
var request_count_thresh int = 1600
var count_mutex sync.Mutex

var all_food_static_string string

func FetchFood(food_id, amount int) (rest_stock int) {
	// return cache.FetchFood(food_id, amount)
	defer func() {
		fmt.Printf("Fetch Food id:%d, amount:%d  res:%d\n", food_id, amount, rest_stock)
	}()
	// Put Back to Cache , not to mem
	if amount < 0 {
		return cache.FetchFood(food_id, amount)
	}

	res := mem.FetchFood(food_id, amount)
	// mem enougth
	if res >= 0 {
		return res
	}
	// mem just fit
	if res == -amount {
		mem.FetchFood(food_id, -amount)
		return cache.FetchFood(food_id, amount)
	}
	// mem not enough, fetch from cache
	mem.FetchFood(food_id, res)
	return cache.FetchFood(food_id, -res)
}

func AllFoods() string {

	if request_count > request_count_thresh {
		return all_food_static_string
	}
	count_mutex.Lock()
	defer count_mutex.Unlock()
	request_count += 1
	if request_count == request_count_thresh {
		spreadFoods()
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
	} else {
		for _, food := range *foods {
			cache.FoodImport(food)
		}
	}

	var generateStaticFoods = func() {
		var fs []common.Food = make([]common.Food, 0)
		for _, food := range *foods {
			food.Price, _ = GetFoodPrice(food.Id)
			food.Stock -= rand.Int() % 512
			fs = append(fs, food)
			mem.SetFoodStock(food.Id, 0)
		}
		ret, _ := json.Marshal(fs)
		all_food_static_string = string(ret)
	}
	generateStaticFoods()

}

func spreadFoods() {
	foods := cache.GetAllFoods()
	var localStockAmount = func(stock int) int {
		switch {
		case stock > 5000:
			return 2000
		case stock > 3000:
			return 1000
		case stock > 1000:
			return 300
		case stock > 500:
			return 200
		case stock > 200:
			return 50
		case stock > 100:
			return 30
			/*
				case stock > 50:
					return 20
				case stock > 20:
					return 5
				case stock > 0:
					return 0
			*/
		}
		return 0
	}

	var spreadFoodStocks = func() {
		for _, food := range *foods {
			am := localStockAmount(food.Stock)
			cache.FetchFood(food.Id, am)
			mem.SetFoodStock(food.Id, am)
		}
	}
	spreadFoodStocks()
}
