package persist

import (
	"common"
	"fmt"
)

func GetAllFoods() *[]common.Food {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()

	db := getDB()
	defer releaseDB(db)

	rows, err := db.Query(s.SELECT_ALL_FOODS)
	if err != nil {
		panic(err)
	}

	var foods []common.Food
	for rows.Next() {
		var food common.Food
		err := rows.Scan(&food.Id, &food.Stock, &food.Price)
		if err != nil {
			fmt.Println(err)
		}
		foods = append(foods, food)
	}
	return &foods
}
