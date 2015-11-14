package server

import (
	"net/http"
)

type Response struct {
	status int
	Code   string `json:"code"`
	Msg    string `json:"message"`
}

var Unauthorized = &Response{
	status: http.StatusUnauthorized,
	Code:   "INVALID_ACCESS_TOKEN",
	Msg:    "无效的令牌",
}

var BadRequest = &Response{
	status: http.StatusBadRequest,
	Code:   "EMPTY_REQUEST",
	Msg:    "请求体为空",
}

var BadFormat = &Response{
	status: http.StatusBadRequest,
	Code:   "MALFORMED_JSON",
	Msg:    "格式错误",
}

var cart_not_exist = &Response{
	status: http.StatusNotFound,
	Code:   "CART_NOT_FOUND",
	Msg:    "篮子不存在",
}

var cart_not_belong = &Response{
	status: http.StatusUnauthorized,
	Code:   "NOT_AUTHORIZED_TO_ACCESS_CART",
	Msg:    "无权限访问指定的篮子",
}

var food_not_enough = &Response{
	status: http.StatusForbidden,
	Code:   "FOOD_OUT_OF_STOCK",
	Msg:    "食物库存不足",
}

var cart_more_than_three = &Response{
	status: http.StatusForbidden,
	Code:   "FOOD_OUT_OF_LIMIT",
	Msg:    "篮子中食物数量超过了三个",
}

var food_not_exist = &Response{
	status: http.StatusNotFound,
	Code:   "FOOD_NOT_FOUND",
	Msg:    "食物不存在",
}

var too_many_orders = &Response{
	status: http.StatusForbidden,
	Code:   "ORDER_OUT_OF_LIMIT",
	Msg:    "每个用户只能下一单",
}

var status map[string]*Response = func() map[string]*Response {
	s := make(map[string]*Response)
	s["cart_not_exist"] = cart_not_exist
	s["cart_not_belong"] = cart_not_belong
	s["cart_more_than_three"] = cart_more_than_three
	s["food_not_enough"] = food_not_enough
	s["food_not_exist"] = food_not_exist
	s["too_many_orders"] = too_many_orders
	return s
}()
