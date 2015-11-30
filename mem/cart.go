package mem

func GetCartRetString(uid int) *string {
	return carts[uid]
}

func SetCartRetString(uid int, str string) {
	carts[uid] = &str
}
