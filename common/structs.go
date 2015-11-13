package common

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
