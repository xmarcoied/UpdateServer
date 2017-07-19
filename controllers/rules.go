package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xmarcoied/go-updater/model"
)

func AddRule(c *gin.Context) {
	c.HTML(http.StatusOK, "newrule.html", gin.H{
		"id": c.Param("id"),
	})

}

func NewRule(c *gin.Context) {
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

	c.Redirect(http.StatusMovedPermanently, "/admin/dashboard/release/"+c.Param("id"))

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
		db.DB.Where("rule_id =?", rule.ID).First(&timerule)
		if !(time.Now().Before(timerule.EndTime) && time.Now().After(timerule.StartTime)) {
			return false
		}
	}
	return true
}
