package core

import (
	"fmt"
	"strconv"

	"code.videolan.org/GSoC2017/Marco/UpdateServer/database"
)

// GetChannels return all channels recorded at the database orded by id
func GetChannels() []database.Channel {
	var channels []database.Channel
	db.DB.Order("id").Find(&channels)
	return channels
}

func GetChannel(name string) database.Channel {
	var channel database.Channel
	db.DB.First(&channel, "name = ?", name)
	return channel
}

// NewChannel create a new channel associated with a public key
func NewChannel(channel *database.Channel) {
	db.DB.Table("channels").Create(&channel)
}

// DeleteChannel
func DeleteChannel(channelName string) {
	// First delete all the releases releases with the 'channelName' channel
	query := fmt.Sprintf("channel = '%s'", channelName)
	releases := GetReleases(query)

	for _, release := range releases {
		DeleteRelease(strconv.Itoa(int(release.ID)))
	}
	// Then delete the channel
	db.DB.Where("name = ?", channelName).Delete(&database.Channel{})
}
