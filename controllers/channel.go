package controllers

import (
	"io/ioutil"
	"net/http"

	"code.videolan.org/GSoC2017/Marco/UpdateServer/model"
	"github.com/gin-gonic/gin"
)

func NewChannelController() *ChannelController {
	return &ChannelController{}
}

func (cc ChannelController) GetChannels(c *gin.Context) {
	var channels []model.Channel
	db.DB.Order("id").Find(&channels)

	c.HTML(http.StatusOK, "channels.html", gin.H{
		"channels": channels,
	})

}

func (cc ChannelController) AddChannel(c *gin.Context) {
	c.HTML(http.StatusOK, "channel.html", nil)

}

func (cc ChannelController) NewChannel(c *gin.Context) {
	var channel model.Channel
	c.Bind(&channel)

	db.DB.Table("channels").Create(&channel)
	pub := "static/channels/public/" + channel.Name + ".asc"
	private := "static/channels/private/" + channel.Name + ".asc"

	ioutil.WriteFile(pub, []byte(channel.PublicKey), 0644)
	ioutil.WriteFile(private, []byte(channel.PrivateKey), 0644)

	c.Redirect(http.StatusMovedPermanently, "/admin/dashboard/channels/")
}
