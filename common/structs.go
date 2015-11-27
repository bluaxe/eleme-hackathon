package common

import (
	"fmt"
	"time"
)

type User struct {
	Id     int
	Name   string
	Passwd string
}

type Food struct {
	Id    int `json:"id"`
	Stock int `json:"stock"`
	Price int `json:"price"`
}

type CartFood struct {
	Id  int `json:"food_id"`
	Num int `json:"count"`
}

type Order struct {
	// Id       int        `json:"-"`
	Id    string     `json:"id"`
	Foods []CartFood `json:"items"`
	Total int        `json:"total"`
}

type AdminOrder struct {
	Order
	Uid int `json:"user_id"`
}

func RecoverAndPrint(msg string) {
	if r := recover(); r != nil {
		// fmt.Println(msg, r)
	}
}

func RecoverPrintDo(msg string, f func()) {
	if r := recover(); r != nil {
		// fmt.Println(msg, r)
		f()
	}
}

func RecoverAndDo(f func()) {
	if r := recover(); r != nil {
		f()
	}
}

func Recover() {
	recover()
}

func LogTime(t time.Time, msg string) {

	now_t := time.Now()
	dur := now_t.Sub(t)
	ms := dur.Nanoseconds() / 1000000
	fmt.Printf("[%dms] %s\n", ms, msg)

}
