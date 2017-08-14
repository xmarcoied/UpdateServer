package main

import (
	"flag"

	"log"
	"math/rand"
	"net/http"
	"time"

	"code.videolan.org/GSoC2017/Marco/UpdateServer/config"
	"code.videolan.org/GSoC2017/Marco/UpdateServer/core"
	"code.videolan.org/GSoC2017/Marco/UpdateServer/database"
	"github.com/gin-gonic/gin"
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

	router := gin.Default()

	router.GET("/admin/dashboard/channels", func(c *gin.Context) {
		channels := core.GetChannels()
		c.HTML(http.StatusOK, "channels.html", gin.H{
			"channels": channels,
		})
	})

	router.GET("/admin/dashboard/channels/add", func(c *gin.Context) {
		c.HTML(http.StatusOK, "channel.html", nil)
	})

	router.POST("/admin/dashboard/new_channel", func(c *gin.Context) {
		var channel database.Channel
		c.Bind(&channel)
		core.NewChannel(channel)
		c.Redirect(http.StatusMovedPermanently, "/admin/dashboard/channels/")
	})

	router.LoadHTMLGlob("html/*.html")
	router.Run(":" + addr)
}
