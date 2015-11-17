package cache

import (
	"common"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"os"
	"time"
)

var addr string

var pool *redis.Pool

func newPool(addr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     10,
		MaxActive:   0,
		IdleTimeout: 10 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func init() {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	addr = fmt.Sprintf("%s:%s", host, port)
	pool = newPool(addr)
}

func SeekMaster() {
	ll := NewLock("seek master")
	if ll.Get() {
		defer ll.Release()
		// fmt.Println("I am master!")
	}
}

func getCon() redis.Conn {
	return pool.Get()
}

func releaseCon(c redis.Conn) {
	c.Close()
}

func Count() int {
	defer common.RecoverAndPrint("Cache Count error")

	c := getCon()
	defer releaseCon(c)

	key := k.COUNT_KEY

	cnt, err := redis.Int(c.Do("incr", key))
	if err != nil {
		panic(err)
	}
	return cnt
}

func GetCount() int {
	defer common.RecoverAndPrint("Cache Get Count Error")
	c := getCon()
	defer releaseCon(c)

	key := k.COUNT_KEY
	cnt, err := redis.Int(c.Do("get", key))
	if err != nil {
		panic(err)
	}
	return cnt
}

func Test() {
	c := getCon()
	defer releaseCon(c)
	fmt.Println("get a conn")
	c.Do("SET", "1", "2")
}
