package cache

import (
	"common"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

func FoodsNum() (n int) {
	/*
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Get Foods Num Err! : ", r)
			}
		}()
	*/
	defer common.Recover()

	c := getCon()
	defer releaseCon(c)

	key := k.FOOD_STOCK_KEY

	n, err := redis.Int(c.Do("hlen", key))
	if err != nil {
		panic(err)
	}
	return n
}

func FoodImport(food common.Food) {
	/*
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Import Food Err! : ", r)
			}
		}()
	*/
	defer common.RecoverAndPrint("Error Import Food Error")

	c := getCon()
	defer releaseCon(c)
	// fmt.Println(food.Id, food.Price, food.Stock)

	key := k.FOOD_STOCK_KEY
	_, err := redis.Int(c.Do("hset", key, food.Id, food.Stock))
	if err != nil {
		panic(err)
	}

	key = k.FOOD_PRICE_KEY
	_, err = redis.Int(c.Do("hset", key, food.Id, food.Price))
	if err != nil {
		panic(err)
	}
}

func GetFoodStock(food_id int) (stock int) {
	defer common.RecoverAndPrint("Not Found Food in Cache")

	c := getCon()
	defer releaseCon(c)

	key := k.FOOD_STOCK_KEY
	stock, err := redis.Int(c.Do("hget", key, food_id))
	if err != nil {
		panic(err)
	}
	return stock
}

func FetchFood(food_id, count int) int {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Cache Fetch Not Found Food in Cache", r)
		}
	}()

	c := getCon()
	defer releaseCon(c)

	key := k.FOOD_STOCK_KEY
	stock, err := redis.Int(c.Do("hincrby", key, food_id, -count))
	fmt.Printf("\t\tCache Fetch fid:%d num:%d, res:%d\n", food_id, count, stock)
	if err != nil {
		panic(err)
	}
	return stock
}

func GetFoodPrice(food_id int) int {
	defer common.RecoverAndPrint("Get Food Price Cache Failed.")
	c := getCon()
	defer releaseCon(c)

	key := k.FOOD_PRICE_KEY

	price, err := redis.Int(c.Do("hget", key, food_id))
	if err != nil {
		panic(err)
	}
	return price
}

func GetAllFoods() *[]common.Food {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Cache Get All Foods Error : ", r)
		}
	}()
	c := getCon()
	defer releaseCon(c)

	key := k.FOOD_STOCK_KEY
	values, err := redis.Values(c.Do("hgetall", key))
	if err != nil {
		panic(err)
	}
	var fstocks []struct {
		Id    int
		Stock int
	}
	if err = redis.ScanSlice(values, &fstocks); err != nil {
		panic(err)
	}

	key = k.FOOD_PRICE_KEY
	values, err = redis.Values(c.Do("hgetall", key))
	if err != nil {
		panic(err)
	}
	var fprice []struct {
		Id    int
		Price int
	}
	if err = redis.ScanSlice(values, &fprice); err != nil {
		panic(err)
	}

	var fs map[int]common.Food = make(map[int]common.Food)
	for _, s := range fstocks {
		f := fs[s.Id]
		f.Stock = s.Stock
		fs[s.Id] = f
	}
	for _, s := range fprice {
		f := fs[s.Id]
		f.Price = s.Price
		fs[s.Id] = f
	}

	var foods []common.Food = make([]common.Food, 0)
	for id, f := range fs {
		f.Id = id
		if f.Stock > 0 {
			foods = append(foods, f)
		}
	}

	return &foods
}

func GetAllFoodsStock() *[]common.Food {
	defer common.RecoverAndPrint("Cache Get All Foods Stock Error!")

	c := getCon()
	defer releaseCon(c)

	key := k.FOOD_STOCK_KEY
	values, err := redis.Values(c.Do("hgetall", key))
	if err != nil {
		panic(err)
	}
	var fstocks []struct {
		Id    int
		Stock int
	}
	if err = redis.ScanSlice(values, &fstocks); err != nil {
		panic(err)
	}

	var foods []common.Food = make([]common.Food, 0)

	for _, s := range fstocks {
		var f common.Food
		f.Id = s.Id
		f.Stock = s.Stock
		foods = append(foods, f)
	}

	return &foods
}
