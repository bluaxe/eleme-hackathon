package service

import (
	"cache"
	"common"
	"encoding/json"
	"fmt"
	"mem"
	"persist"
	"time"
)

var food_list = func() map[int]int {
	var f = make(map[int]int)
	return f
}()

var all_food_static_string string

func FetchFood(food_id, amount int) (rest_stock int) {
	// return cache.FetchFood(food_id, amount)

	var spread = func(cache_res int) {
		if localStockAmount(cache_res) > 0 {
			// fmt.Printf("spread food from cache to mem id:%d mount:%d\n", food_id, localStockAmount(cache_res))
			var food = &common.Food{
				Id:    food_id,
				Stock: cache_res,
			}
			spreadFoodStock(food)
		}
	}

	// defer common.RecoverPrintDo("Error Fetch Food in service. ", func() { rest_stock = 0 })

	// Put Back to Cache , not to mem
	if amount < 0 {
		return cache.FetchFood(food_id, amount)
		// fmt.Printf("Fetch Food id:%d, amount:%d  Put Back To Cache\n", food_id, amount)
	}

	got := mem.FetchFoodToEmpty(food_id, amount)
	// mem enougth
	if got == amount {
		// fmt.Printf("Fetch Food id:%d, amount:%d  got:%d From Mem\n", food_id, amount, got)
		return 1
	}

	// mem just fit
	if got == 0 {
		// fmt.Printf("Fetch Food id:%d, amount:%d  mem is 0, get from cache\n", food_id, amount)
		res := cache.FetchFood(food_id, amount)
		go spread(res)
		return res
	}

	if got < amount {
		// mem not enough, fetch from cache
		// fmt.Printf("Fetch Food id:%d, amount:%d  mem is not enough, get %d from cache\n", food_id, amount, amount-got)
		res := cache.FetchFood(food_id, amount-got)
		go spread(res)
		return res
	}
	// panic(fmt.Sprintf("Unknown error !!! food_id:%d, amount:%d, mem got : %d", food_id, amount, got))
	return 0
}

func AllFoods() string {
	// return static data since this interface is not important.
	return all_food_static_string

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
			mem.AddFoodStock(food.Id, 0)
		}
		ret, _ := json.Marshal(fs)
		all_food_static_string = string(ret)
	}
	generateStaticFoods()
	spreadFoods()
	go RefreshFoodList()
}

func RefreshFoodList() {
	fmt.Println("Start time:", time.Now())
	ticker := time.NewTicker(10 * time.Second)
	for _ = range ticker.C {
		foods := cache.GetAllFoodsStock()
		var fs []common.Food = make([]common.Food, 0)
		for _, food := range *foods {
			food.Price, _ = GetFoodPrice(food.Id)
			/*
				add estimated stock in memory
			*/
			fs = append(fs, food)
		}
		ret, _ := json.Marshal(fs)
		all_food_static_string = string(ret)
	}
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
		// case stock > 50:
		// return 10
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
	mem.AddFoodStock(food.Id, am)
}

func spreadFoods() {
	foods := cache.GetAllFoods()
	for _, food := range *foods {
		spreadFoodStock(&food)
	}
}
