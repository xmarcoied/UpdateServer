package core

import (
	"fmt"

	"code.videolan.org/GSoC2017/Marco/UpdateServer/database"
	"code.videolan.org/GSoC2017/Marco/UpdateServer/utils"
)

// GetChannels return all channels recorded at the database orded by id
func GetChannels() []database.Channel {
	var channels []database.Channel
	db.DB.Order("id").Find(&channels)
	return channels
}

// GetChannel return the channel database struct
func GetChannel(name string) database.Channel {
	var channel database.Channel
	db.DB.First(&channel, "name = ?", name)
	return channel
}

// NewChannel create a new channel associated with a public key
func NewChannel(channel *database.Channel) {
	db.DB.Table("channels").Create(&channel)
}

// DeleteChannel delete the channel given by its name
func DeleteChannel(channelName string) {
	// First delete all the releases releases with the 'channelName' channel
	query := fmt.Sprintf("channel = '%s'", channelName)
	releases := GetReleases(query)

	for _, release := range releases {
		DeleteRelease(release.IdToString())
	}
	// Then delete the channel
	db.DB.Where("name = ?", channelName).Delete(&database.Channel{})
}

// CheckChannel Check the channel validity in terms of gpg and uniqueness of its name
func CheckChannel(channel database.Channel) (bool, error) {
	var count int
	db.DB.Model(&database.Channel{}).Where("name = ?", channel.Name).Count(&count)
	if count != 0 {
		return false, fmt.Errorf("Can't have a duplicated channel name")
	} else if err := utils.CheckPGPKey(channel.PublicKey); err != nil {
		return false, fmt.Errorf("public key error: %v", err)
	} else if err := utils.CheckPGPKey(channel.PrivateKey); err != nil && channel.PrivateKey != "" {
		return false, fmt.Errorf("private key error: %v", err)
	} else {
		return true, nil
	}
}
