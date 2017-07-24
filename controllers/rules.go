package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xmarcoied/go-updater/model"
)

func NewRulesController() *RulesController {
	return &RulesController{}
}

func (rsc RulesController) AddRule(c *gin.Context) {
	c.HTML(http.StatusOK, "newrule.html", gin.H{
		"id": c.Param("id"),
	})

}

func (rsc RulesController) NewRule(c *gin.Context) {

	switch c.Param("rule") {
	case "time":
		var buf struct {
			Start string `form:"timestart"`
			End   string `form:"timeend"`
		}
		c.Bind(&buf)

		layout := "2006-01-02T15:04"
		t_start, _ := time.Parse(layout, buf.Start)
		t_end, _ := time.Parse(layout, buf.End)

		var release model.Release
		db.DB.Where("id = ?", c.Param("id")).First(&release)

		release.Rules.TimeRule.StartTime = t_start

		if t_end.IsZero() {
			t_end, _ = time.Parse(layout, "2906-01-02T15:04")
		}
		release.Rules.TimeRule.EndTime = t_end
		db.DB.Save(&release)

	case "os":
		var buf struct {
			OsVersion string `form:"os_version"`
		}
		c.Bind(&buf)

		var release model.Release
		db.DB.Where("id = ?", c.Param("id")).First(&release)
		release.Rules.OsRule.OsVersion = buf.OsVersion

		db.DB.Save(&release)

	case "version":
		var buf struct {
			ProductVersion string `form:"product_version"`
		}
		c.Bind(&buf)

		var release model.Release
		db.DB.Where("id = ?", c.Param("id")).First(&release)
		release.Rules.VersionRule.ProductVersion = buf.ProductVersion

		db.DB.Save(&release)
	case "ip":
		var buf struct {
			IP string `form:"ip_address"`
		}
		c.Bind(&buf)

		var release model.Release
		db.DB.Where("id = ?", c.Param("id")).First(&release)
		release.Rules.IPRule.IP = buf.IP

		db.DB.Save(&release)
	}

	c.Redirect(http.StatusMovedPermanently, "/admin/dashboard/release/"+c.Param("id"))

}

func (rsc RulesController) DeleteRule(c *gin.Context) {
	switch c.Param("rule") {
	case "time":
		var rule model.TimeRule
		db.DB.Where("rule_id = ?", c.Param("id")).Delete(&rule)
	case "os":
		var rule model.OsRule
		db.DB.Where("rule_id = ?", c.Param("id")).Delete(&rule)
	case "version":
		var rule model.VersionRule
		db.DB.Where("rule_id = ?", c.Param("id")).Delete(&rule)
	case "ip":
		var rule model.IPRule
		db.DB.Where("rule_id = ?", c.Param("id")).Delete(&rule)
	}
	c.Redirect(http.StatusMovedPermanently, "/admin/dashboard/releases")

}

func CountRules(release model.Release) int {
	var ret int
	db.DB.Model(&model.Rule{}).Where("release_id = ?", release.ID).Count(&ret)
	return ret
}

func CheckTimeRule(release model.Release) bool {
	var rules []model.Rule
	var timerule model.TimeRule
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

func CheckVersionRule(release model.Release) bool {
	var rules []model.Rule
	var versionrule model.VersionRule
	db.DB.Where("release_id = ?", release.ID).Find(&rules)
	for _, rule := range rules {
		if err := db.DB.Where("rule_id =?", rule.ID).First(&versionrule).Error; err == nil {
			if release.ProductVersion == versionrule.ProductVersion {
				return false
			}
		}
	}

	return true
}

func CheckOsRule(release model.Release) bool {
	var rules []model.Rule
	var osrule model.OsRule
	db.DB.Where("release_id = ?", release.ID).Find(&rules)
	for _, rule := range rules {
		if err := db.DB.Where("rule_id =?", rule.ID).First(&osrule).Error; err == nil {
			if release.OsVer == osrule.OsVersion {
				return false
			}
		}
	}

	return true
}

func CheckIPRule(release model.Release, request model.UpdateRequest) (bool, bool) {
	var rules []model.Rule
	var iprule model.IPRule

	found := false
	db.DB.Where("release_id = ?", release.ID).Find(&rules)
	for _, rule := range rules {
		if err := db.DB.Where("rule_id =?", rule.ID).First(&iprule).Error; err == nil {
			found = true
			if request.IP == iprule.IP {
				return true, true
			}
		}
	}

	return found, false
}
