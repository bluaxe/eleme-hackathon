package cache

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"os"
)

var addr string

func init() {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	addr = fmt.Sprintf("%s:%s", host, port)
}

func SeekMaster() {
	ll := NewLock("seek master")
	if ll.Get() {
		defer ll.Release()
		fmt.Println("I am master!")
	}
}

func getCon() redis.Conn {
	c, err := redis.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	return c
}

func releaseCon(c redis.Conn) {
	c.Close()
}

func Test() {
	c := getCon()
	defer releaseCon(c)
	fmt.Println("get a conn")
	c.Do("SET", "1", "2")
}
