package mem

func InitOrderSlice(size int) {
	orders = make(map[int]*string)
}

func GetOrderRetString(uid int) *string {
	return orders[uid]
}

func SetOrderRetString(uid int, str string) {
	orders[uid] = &str
}
