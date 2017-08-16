package core

import (
	"log"
	"math/rand"
	"time"

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

func CheckTimeRule(release database.Release) bool {
	var rules []database.Rule
	var timerule database.TimeRule
	db.DB.Where("release_id = ?", release.ID).Find(&rules)
	for _, rule := range rules {
		if err := db.DB.Where("rule_id =?", rule.ID).First(&timerule).Error; err == nil {
			if !(time.Now().Before(timerule.EndTime) && time.Now().After(timerule.StartTime)) {
				return false
			}
		}
	}
	return true
}

func CheckOsRule(release database.Release, request database.UpdateRequest) bool {
	var rules []database.Rule
	var osrule database.OsRule
	db.DB.Where("release_id = ?", release.ID).Find(&rules)
	for _, rule := range rules {
		if err := db.DB.Where("rule_id =?", rule.ID).First(&osrule).Error; err == nil {
			log.Println("hola os rule : ", osrule.OsVersion, "hola os request:", release.OsVer, "Compression:", release.OsVer == osrule.OsVersion)
			if request.OsVer == osrule.OsVersion {
				log.Println("That suppose to be a false response")
				return false
			}
		}
	}

	return true
}

func CheckRollRule(release database.Release) bool {
	var rules []database.Rule
	var rollrule database.RollRule

	db.DB.Where("release_id = ?", release.ID).Find(&rules)
	for _, rule := range rules {
		if err := db.DB.Where("rule_id =?", rule.ID).First(&rollrule).Error; err == nil {
			if rand.Intn(100) > rollrule.RollingPercentage {
				return false
			}
		}
	}

	return true
}

func CheckVersionRule(release database.Release, request database.UpdateRequest) bool {
	var rules []database.Rule
	var versionrule database.VersionRule
	db.DB.Where("release_id = ?", release.ID).Find(&rules)
	for _, rule := range rules {
		if err := db.DB.Where("rule_id =?", rule.ID).First(&versionrule).Error; err == nil {
			if request.ProductVersion == versionrule.ProductVersion {
				return false
			}
		}
	}

	return true
}
