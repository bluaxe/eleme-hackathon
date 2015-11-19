package mem

import (
	"common"
	"sync"
)

func FetchFood(food_id, amount int) (res int) {
	defer common.RecoverPrintDo("mem fetch food error.", func() { res = -1 })
	food_lock[food_id].Lock()
	defer food_lock[food_id].Unlock()
	food_stock[food_id] -= amount
	return food_stock[food_id]
}

func SetFoodStock(food_id, stock int) {
	_, have := food_stock[food_id]
	food_stock[food_id] = stock
	if have {
		return
	}
	food_lock[food_id] = &sync.Mutex{}
	food_lock[food_id].Lock()
	food_lock[food_id].Unlock()
}
