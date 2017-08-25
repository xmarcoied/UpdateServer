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
	fingerprint, err := utils.GetFingerprint(channel)
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
	core.NewChannel(&channel)
	c.Redirect(http.StatusMovedPermanently, "/admin/dashboard/channel/"+channel.Name)
}

// DeleteChannel
func DeleteChannel(c *gin.Context) {
	core.DeleteChannel(c.Param("name"))
}

// CheckChannel
func CheckChannel(c *gin.Context) {
	var channel database.Channel
	c.Bind(&channel)
	ret, err := core.CheckChannel(channel)
	if ret == true {
		c.String(http.StatusOK, "OK")
	} else {
		c.String(http.StatusOK, err.Error())
	}

}
