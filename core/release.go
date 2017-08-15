package core

import (
	"code.videolan.org/GSoC2017/Marco/UpdateServer/database"
)

// GetReleases return all releases recorded at the database orded by id
func GetReleases() []database.Release {
	var releases []database.Release
	db.DB.Order("id").Find(&releases)
	return releases
}
