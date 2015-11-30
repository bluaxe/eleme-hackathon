package main

import (
	// "cache"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"server"
	"service"
)

type initializer func()

var jobs map[string]initializer = make(map[string]initializer)

func init() {
	// jobs["cache_seek_master"] = cache.SeekMaster
	jobs["cache_load_foods"] = service.InitFoodsFromPersist
	jobs["mem_load_users"] = service.LoadAllUserToMem
	jobs["mem_gen_orders"] = server.InitOrderStrings
	jobs["mem_gen_login_rets"] = server.InitLoginRetStrings
	jobs["mem_gen_cart_rets"] = server.InitCartRetStrings
	// jobs["mem_gen_tokens"] = service.GenerateTokens
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	host := os.Getenv("APP_HOST")
	port := os.Getenv("APP_PORT")
	debug.SetGCPercent(500)
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf("%s:%s", host, port)

	service.GenerateTokens()
	for _, v := range jobs {
		v()
	}

	server.Start(addr)
}
