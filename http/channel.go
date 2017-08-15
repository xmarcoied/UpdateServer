package http

import (
	"net/http"

	"code.videolan.org/GSoC2017/Marco/UpdateServer/core"
	"code.videolan.org/GSoC2017/Marco/UpdateServer/database"
	"github.com/gin-gonic/gin"
)

// GetChannels is http handler to represent all the channels available in the UpdateServer
func GetChannels(c *gin.Context) {
	channels := core.GetChannels()
	c.HTML(http.StatusOK, "channels.html", gin.H{
		"channels": channels,
	})
}

// AddChannel is http handler to show "New Channel" html page
func AddChannel(c *gin.Context) {
	c.HTML(http.StatusOK, "channel.html", nil)
}

// NewChannel is http handler binding the channel data to create a new channel
func NewChannel(c *gin.Context) {
	var channel database.Channel
	c.Bind(&channel)
	core.NewChannel(channel)
	c.Redirect(http.StatusMovedPermanently, "/admin/dashboard/channels/")
}
