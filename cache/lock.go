package cache

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

var exp = 20000 // int px (3 secondes)
var lock_wait_interval = 1 * time.Millisecond

type Lock struct {
	Key    string
	value  string
	expire int
}

func NewLock(key string) *Lock {
	var l Lock = Lock{
		Key:    key,
		expire: exp,
	}
	return &l
}

func NewLockWithExpire(key string, input_exp int) *Lock {
	var l Lock = Lock{
		Key:    key,
		expire: input_exp,
	}
	return &l
}

func (l *Lock) GetWait() {
	for {
		if l.Get() {
			break
		}
		time.Sleep(lock_wait_interval)
	}
}

func (l *Lock) Get() (ok bool) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Panic in Lock !!!  : ", r)
		}
	}()

	con := getCon()
	defer releaseCon(con)

	lockkey := getLockKey(l.Key)

	rand := getRandString()

	ret, err := con.Do("SET", lockkey, rand, "NX", "PX", l.expire)
	if err != nil {
		panic(err)
	}
	_, err = redis.String(ret, err)
	if err != redis.ErrNil {
		l.value = rand
		// fmt.Printf("[%s]Got Lock for key: %s \n", time.Now().Format(time.UnixDate), l.Key)
		return true
	} else {
		return false
	}
}

func (l *Lock) Release() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Panic in Lock Release!", r)
		}
	}()
	c := getCon()
	defer releaseCon(c)

	lockkey := getLockKey(l.Key)
	ret, err := c.Do("eval", k.RELEASE_LOCK, 1, lockkey, l.value)
	if err != nil {
		panic(err)
	}
	_, err = redis.Int(ret, err)
	if err != nil {
		panic(err)
	}
	// fmt.Printf("Release Lock for key: %s result: %d\n", l.Key, ok)
}
