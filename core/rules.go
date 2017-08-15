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

func GetRules(release database.Release) []database.Rule {
	var rules []database.Rule
	db.DB.Model(&release).Related(&rules)
	for i, rule := range rules {
		db.DB.Model(&rule).Related(&rules[i].TimeRule)
		db.DB.Model(&rule).Related(&rules[i].RollRule)
		db.DB.Model(&rule).Related(&rules[i].VersionRule)
		db.DB.Model(&rule).Related(&rules[i].OsRule)
		db.DB.Model(&rule).Related(&rules[i].IPRule)
	}
	return rules
}

func DeleteRule(rule string, id string) {
	switch rule {
	case "time":
		var rule database.TimeRule
		db.DB.Where("rule_id = ?", id).Delete(&rule)
	case "os":
		var rule database.OsRule
		db.DB.Where("rule_id = ?", id).Delete(&rule)
	case "version":
		var rule database.VersionRule
		db.DB.Where("rule_id = ?", id).Delete(&rule)
	case "ip":
		var rule database.IPRule
		db.DB.Where("rule_id = ?", id).Delete(&rule)
	case "roll":
		var rule database.RollRule
		db.DB.Where("rule_id = ?", id).Delete(&rule)
	}
}
