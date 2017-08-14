package http

import (
	"net/http"

	"code.videolan.org/GSoC2017/Marco/UpdateServer/core"
	"code.videolan.org/GSoC2017/Marco/UpdateServer/database"
	"github.com/gin-gonic/gin"
)

func RouterInit() *gin.Engine {
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
	return router
}

func Run(addr string) {
	RouterInit().Run(":" + addr)
}
