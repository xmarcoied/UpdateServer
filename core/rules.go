package core

import (
	"code.videolan.org/GSoC2017/Marco/UpdateServer/database"
)

// NewRule function create a new rule assoicated with a release. Given the release_id
func NewRule(release_id string, rule database.Rule) {
	var release database.Release
	db.DB.Where("id = ?", release_id).First(&release)

	release.Rules = append(release.Rules, rule)
	db.DB.Save(&release)
}

// GetRules function return all the rules associated with a release. Given the release
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

// DeleteRule function delete a specific rule from a release. Given the rule type and rule_id
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
