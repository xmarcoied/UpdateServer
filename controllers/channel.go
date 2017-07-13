package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xmarcoied/go-updater/model"
)

func GetChannels(c *gin.Context) {
	var channels []model.Channel
	db.DB.Order("id").Find(&channels)

	c.HTML(http.StatusOK, "channels.html", gin.H{
		"channels": channels,
	})

}

func AddChannel(c *gin.Context) {
	c.HTML(http.StatusOK, "channel.html", nil)

}

func NewChannel(c *gin.Context) {
	var channel model.Channel
	c.Bind(&channel)
	log.Println(channel)

	db.DB.Table("channels").Create(&channel)
	c.Redirect(http.StatusMovedPermanently, "/admin/dashboard/channels/")
}
