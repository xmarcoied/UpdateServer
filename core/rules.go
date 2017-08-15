package core

import (
	"code.videolan.org/GSoC2017/Marco/UpdateServer/database"
)

func NewRule(release_id string, rule database.Rule) {
	var release database.Release
	db.DB.Where("id = ?", release_id).First(&release)

	release.Rules = append(release.Rules, rule)
	db.DB.Save(&release)
}
