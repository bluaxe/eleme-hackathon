package mem

import (
	"common"
	"fmt"
	"sync"
)

func FetchFoodToEmpty(food_id, amount int) (got int) {
	defer common.RecoverPrintDo("mem fetch to empty food error.", func() { got = 0 })
	food_lock[food_id].Lock()
	defer food_lock[food_id].Unlock()
	got = 0
	if food_stock[food_id] == 0 {
		return 0
	}
	if food_stock[food_id] < 0 {
		panic(fmt.Sprintf("Error !!!!! food stock of %d is negetive : %d\n", food_id, food_stock[food_id]))
	}
	if food_stock[food_id] < amount {
		got = food_stock[food_id]
		food_stock[food_id] = 0
	} else {
		got = amount
		food_stock[food_id] -= amount
	}
	return got
}

func AddFoodStock(food_id, stock int) {
	_, have := food_stock[food_id]
	if !have {
		food_lock[food_id] = &sync.Mutex{}
		food_stock[food_id] = 0
	}
	food_lock[food_id].Lock()
	defer food_lock[food_id].Unlock()
	food_stock[food_id] += stock
}
