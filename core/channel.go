package core

import (
	"io/ioutil"

	"code.videolan.org/GSoC2017/Marco/UpdateServer/database"
)

func GetChannels() []database.Channel {
	var channels []database.Channel
	db.DB.Order("id").Find(&channels)
	return channels
}

func NewChannel(channel database.Channel) {
	db.DB.Table("channels").Create(&channel)
	pub := "static/channels/public/" + channel.Name + ".asc"

	ioutil.WriteFile(pub, []byte(channel.PublicKey), 0644)
}
