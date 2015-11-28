package mem

import (
	"common"
	"fmt"
	"sync"
	"sync/atomic"
)

func FetchFoodToEmpty(food_id, am int) (got int) {
	amount := int32(am)
	defer common.RecoverPrintDo("mem fetch to empty food error.", func() { got = 0 })
	if *food_stock[food_id] > int32(30) {
		atomic.AddInt32(food_stock[food_id], -amount)
		return am
	}
	food_lock[food_id].Lock()
	defer food_lock[food_id].Unlock()
	got = 0
	if *food_stock[food_id] == int32(0) {
		return 0
	}
	if *food_stock[food_id] < int32(0) {
		panic(fmt.Sprintf("Error !!!!! food stock of %d is negetive : %d\n", food_id, *food_stock[food_id]))
	}
	if *food_stock[food_id] < amount {
		got = int(*food_stock[food_id])
		*food_stock[food_id] = 0
	} else {
		got = int(amount)
		*food_stock[food_id] -= amount
	}
	return got
}

func AddFoodStock(food_id, stock int) {
	_, have := food_stock[food_id]
	if !have {
		food_lock[food_id] = &sync.Mutex{}
		var stock int32 = 0
		food_stock[food_id] = &stock
	}
	food_lock[food_id].Lock()
	defer food_lock[food_id].Unlock()
	*food_stock[food_id] += int32(stock)
}
