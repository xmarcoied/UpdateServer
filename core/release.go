package core

import (
	"code.videolan.org/GSoC2017/Marco/UpdateServer/database"
)

func GetReleases() []database.Release {
	var releases []database.Release
	db.DB.Order("id").Find(&releases)
	return releases
}
