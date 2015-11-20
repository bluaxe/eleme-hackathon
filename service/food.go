package service

import (
	"cache"
	"common"
	"encoding/json"
	"fmt"
	// "math/rand"
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
	/*
		defer func() {
			fmt.Printf("Fetch Food id:%d, amount:%d  res:%d\n", food_id, amount, rest_stock)
		}()
	*/

	// Put Back to Cache , not to mem
	if amount < 0 {
		return cache.FetchFood(food_id, amount)
		fmt.Printf("Fetch Food id:%d, amount:%d  Put Back To Cache\n", food_id, amount)
	}

	res := mem.FetchFood(food_id, amount)
	// mem enougth
	if res >= 0 {
		fmt.Printf("Fetch Food id:%d, amount:%d  res:%d In Mem\n", food_id, amount, res)
		return res
	}
	// mem just fit
	if res == -amount {
		mem.FetchFood(food_id, -amount)
		fmt.Printf("Fetch Food id:%d, amount:%d  mem is 0, get from cache\n", food_id, amount)
		return cache.FetchFood(food_id, amount)
	}
	// mem not enough, fetch from cache
	fmt.Printf("Fetch Food id:%d, amount:%d  mem is not enough, get %d from cache\n", food_id, amount, -res)
	mem.FetchFood(food_id, res)
	cache_res := cache.FetchFood(food_id, -res)

	if localStockAmount(cache_res) > 0 {
		fmt.Printf("spread food from cache to mem id:%d mount:%d\n", food_id, localStockAmount(cache_res))
		var food = &common.Food{
			Id:    food_id,
			Stock: cache_res,
		}
		spreadFoodStock(food)
	}
	return cache_res
}

func AllFoods() string {
	// return static data since this interface is not important.
	return all_food_static_string

	// count request and when > thresh , return static data
	/*
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
	*/

	// real time data from cache
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
			fs = append(fs, food)
			mem.SetFoodStock(food.Id, 0)
		}
		ret, _ := json.Marshal(fs)
		all_food_static_string = string(ret)
	}
	generateStaticFoods()
	spreadFoods()
}

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
		return 40
	case stock > 100:
		return 20
	case stock > 50:
		return 10
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

var spreadFoodStock = func(food *common.Food) {
	am := localStockAmount(food.Stock)
	cache.FetchFood(food.Id, am)
	mem.SetFoodStock(food.Id, am)
}

func spreadFoods() {
	foods := cache.GetAllFoods()
	for _, food := range *foods {
		spreadFoodStock(&food)
	}
}
