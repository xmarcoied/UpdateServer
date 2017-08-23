package http

import (
	"log"
	"net/http"

	"code.videolan.org/GSoC2017/Marco/UpdateServer/core"
	"code.videolan.org/GSoC2017/Marco/UpdateServer/database"
	"code.videolan.org/GSoC2017/Marco/UpdateServer/utils"
	"github.com/gin-gonic/gin"
)

// GetChannels is http handler to represent all the channels available in the UpdateServer
func GetChannels(c *gin.Context) {
	channels := core.GetChannels()
	c.HTML(http.StatusOK, "channels.html", gin.H{
		"channels": channels,
	})
}

// GetChannel is http handler to represent channel content
func GetChannel(c *gin.Context) {
	channel := core.GetChannel(c.Param("name"))
	fingerprint, err := utils.GetFingerprint(c.Param("name"))
	if err != nil {
		log.Println(err)
	}
	c.HTML(http.StatusOK, "channel.html", gin.H{
		"title":       "view channel",
		"channel":     channel,
		"fingerprint": fingerprint,
	})
}

// AddChannel is http handler to show "New Channel" html page
func AddChannel(c *gin.Context) {
	c.HTML(http.StatusOK, "channel.html", gin.H{
		"title": "new channel",
	})
}

// NewChannel is http handler binding the channel data to create a new channel
func NewChannel(c *gin.Context) {
	var channel database.Channel
	c.Bind(&channel)
	core.NewChannel(channel)
	c.Redirect(http.StatusMovedPermanently, "/admin/dashboard/channels/")
}

// DeleteChannel
func DeleteChannel(c *gin.Context) {
	core.DeleteChannel(c.Query("name"))
	c.Redirect(http.StatusMovedPermanently, "/admin/dashboard/channels")
}
