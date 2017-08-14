package http

import (
	"net/http"

	"code.videolan.org/GSoC2017/Marco/UpdateServer/core"
	"code.videolan.org/GSoC2017/Marco/UpdateServer/database"
	"github.com/gin-gonic/gin"
)

func GetChannels(c *gin.Context) {
	channels := core.GetChannels()
	c.HTML(http.StatusOK, "channels.html", gin.H{
		"channels": channels,
	})
}

func AddChannel(c *gin.Context) {
	c.HTML(http.StatusOK, "channel.html", nil)
}

func NewChannel(c *gin.Context) {
	var channel database.Channel
	c.Bind(&channel)
	core.NewChannel(channel)
	c.Redirect(http.StatusMovedPermanently, "/admin/dashboard/channels/")
}
