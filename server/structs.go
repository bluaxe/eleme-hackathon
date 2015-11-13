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
