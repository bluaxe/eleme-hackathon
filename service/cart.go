package service

import (
	"cache"
	"common"
	"crypto/rand"
	"fmt"
	"strings"
)

func saveCart(uid int, cid string) {
	cache.SaveCart(uid, cid)
}

func NewCart(uid int) string {
	b := make([]byte, 2)
	rand.Read(b)
	cart_id := fmt.Sprintf("%dsjtu%x", uid*2+1, b)
	// saveCart(id, cart_id)
	return cart_id
}

func DestroyCart(cart_id string) {
	cache.DelCart(cart_id)
}

func GetCartUser(cart_id string) (id int) {
	defer common.RecoverPrintDo("Service get cart user id failed", func() { id = -1 })
	id_str := strings.Split(cart_id, "sjtu")[0]
	fmt.Sscanf(id_str, "%d", &id)
	return (id - 1) / 2
	// return cache.GetCartUser(cart_id)
}

func AddFood(fid, count, uid int, cid string) string {
	cart_uid := GetCartUser(cid)
	if cart_uid == -1 {
		return "cart_not_exist"
	}
	if cart_uid != uid {
		return "cart_not_belong"
	}

	_, exist := food_list[fid]
	if !exist {
		return "food_not_exist"
	}

	foods := cache.GetCartFoods(cid)
	// fmt.Println("cart foods len", len(*foods))
	sum := 0
	for _, f := range *foods {
		sum += f.Num
	}

	if sum+count > 3 {
		return "cart_more_than_three"
	}

	cache.CartAddFood(cid, fid, count)

	return "ok"
}
