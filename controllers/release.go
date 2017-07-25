package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xmarcoied/go-updater/model"
	"github.com/xmarcoied/go-updater/utils"
)

func NewReleaseController() *ReleaseController {
	return &ReleaseController{}
}

// Show all releases
func (rlc ReleaseController) GetReleases(c *gin.Context) {
	var releases []model.Release
	db.DB.Order("id").Find(&releases)

	c.HTML(http.StatusOK, "releases.html", gin.H{
		"releases": releases,
	})
}

func (rlc ReleaseController) GetRelease(c *gin.Context) {
	var (
		release     []model.Release
		rules       []model.Rule
		timerule    []model.TimeRule
		osrule      []model.OsRule
		versionrule []model.VersionRule
		iprule      []model.IPRule
		rollrule    []model.RollRule
		channels    []model.Channel
	)
	if err := db.DB.Where("id = ?", c.Param("id")).Find(&release).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		log.Println(err)

	} else {
		// FIXME : Must be an implementation better than this.
		var buf model.TimeRule
		db.DB.Where("release_id = ?", c.Param("id")).Find(&rules)
		for _, rule := range rules {
			if err := db.DB.Model(buf).Where("rule_id =?", rule.ID).First(&buf).Error; err == nil {
				timerule = append(timerule, buf)
			}
		}

		var buf2 model.OsRule
		for _, rule := range rules {
			if err := db.DB.Model(buf2).Where("rule_id =?", rule.ID).First(&buf2).Error; err == nil {
				osrule = append(osrule, buf2)
			}
		}

		var buf3 model.VersionRule
		for _, rule := range rules {
			if err := db.DB.Model(buf3).Where("rule_id =?", rule.ID).First(&buf3).Error; err == nil {
				versionrule = append(versionrule, buf3)
			}
		}

		var buf4 model.IPRule
		for _, rule := range rules {
			if err := db.DB.Model(buf4).Where("rule_id =?", rule.ID).First(&buf4).Error; err == nil {
				iprule = append(iprule, buf4)
			}
		}

		var buf5 model.RollRule
		for _, rule := range rules {
			if err := db.DB.Model(buf5).Where("rule_id =?", rule.ID).First(&buf5).Error; err == nil {
				rollrule = append(rollrule, buf5)
			}
		}

		db.DB.Model(&channels).Find(&channels)

		c.HTML(http.StatusOK, "release.html", gin.H{
			"id":          c.Param("id"),
			"release":     release,
			"timerule":    timerule,
			"osrule":      osrule,
			"versionrule": versionrule,
			"iprule":      iprule,
			"rollrule":    rollrule,
			"channels":    channels,
		})
	}

}

func (rlc ReleaseController) EditRelease(c *gin.Context) {
	var release model.Release
	c.Bind(&release)

	db.DB.Model(&release).Where("id = ?", c.Param("id")).Updates(release)

	id_buf, _ := strconv.Atoi(c.Param("id"))
	release.ID = uint(id_buf)

	GenerateStatus(release)
	GenerateSignature(release)
	c.Redirect(http.StatusMovedPermanently, "/admin/dashboard/releases/")

}

func (rlc ReleaseController) DelRelease(c *gin.Context) {
	var release model.Release

	db.DB.Where("id = ?", c.Param("id")).Delete(&release)
	c.Redirect(http.StatusMovedPermanently, "/admin/dashboard/releases/")
}

// New release
func (rlc ReleaseController) NewRelease(c *gin.Context) {
	var release model.Release
	c.Bind(&release)

	db.DB.Save(&release)
	GenerateStatus(release)
	GenerateSignature(release)
	if utils.ProcessRelease(release) == true {
		log.Println("Accepted Release")
		c.Redirect(http.StatusMovedPermanently, "/admin/dashboard/releases/")

	} else {
		log.Println("Refused Release")
		c.Redirect(http.StatusMovedPermanently, "/admin/dashboard/newrelease")
	}
}

// Admin dashboard (new releases)
func (rlc ReleaseController) AddRelease(c *gin.Context) {
	var channels []model.Channel
	db.DB.Model(&channels).Find(&channels)

	c.HTML(http.StatusOK, "newrelease.html", gin.H{
		"channels": channels,
	})
}

func GenerateStatus(release model.Release) {
	release.Signature = ""
	ReleaseDir := "static/releases/" + strconv.Itoa(int(release.ID))
	ReleaseJSON, _ := json.Marshal(release)
	ioutil.WriteFile(ReleaseDir, ReleaseJSON, 0644)

}

func GenerateSignature(release model.Release) {
	SignatureDir := "static/signatures/" + strconv.Itoa(int(release.ID)) + ".asc"
	ReleaseDir := "static/releases/" + strconv.Itoa(int(release.ID))
	PrivateKeyDir := "static/channels/private/" + release.Channel + ".asc"

	err := utils.Sign(PrivateKeyDir, ReleaseDir, SignatureDir)
	log.Println(err)
}
