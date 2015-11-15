package service

import (
	"cache"
	"crypto/rand"
	"fmt"
)

func saveCart(uid int, cid string) {
	cache.SaveCart(uid, cid)
}

func NewCart(id int) string {
	b := make([]byte, 6)
	rand.Read(b)
	cart_id := fmt.Sprintf("%x", b)
	go saveCart(id, cart_id)
	return cart_id
}

func DestroyCart(cart_id string) {
	cache.DelCart(cart_id)
}

func AddFood(fid, count, uid int, cid string) string {
	cart_uid := cache.GetCartUser(cid)
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

	if count > 3 {
		return "cart_more_than_three"
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

	go cache.CartAddFood(cid, fid, count)

	return "ok"
}
