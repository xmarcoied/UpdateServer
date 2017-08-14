package main

import (
	"flag"
	"log"
	"math/rand"
	"time"

	"code.videolan.org/GSoC2017/Marco/UpdateServer/config"
	"code.videolan.org/GSoC2017/Marco/UpdateServer/core"
	"code.videolan.org/GSoC2017/Marco/UpdateServer/database"
	"code.videolan.org/GSoC2017/Marco/UpdateServer/http"
)

var (
	db   database.Impl
	addr string
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	flag.StringVar(&addr, "port", "8080", "The port server will be running on")
	flag.Parse()
}

func main() {
	c, err := config.Get()
	if err != nil {
		log.Fatal(err)
	}

	err = db.ConnectDB(c)
	if err != nil {
		log.Fatal(err)
	}

	core.SetDB(db)
	http.Run(addr)
}
