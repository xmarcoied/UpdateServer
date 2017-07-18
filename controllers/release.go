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

// Show all releases
func GetReleases(c *gin.Context) {
	var releases []model.Release
	db.DB.Order("id").Find(&releases)

	c.HTML(http.StatusOK, "releases.html", gin.H{
		"releases": releases,
	})
}

func GetRelease(c *gin.Context) {
	var (
		release  []model.Release
		rules    []model.Rule
		timerule []model.TimeRule
		buf      model.TimeRule
		RulesID  uint
	)
	if err := db.DB.Where("id = ?", c.Param("id")).Find(&release).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		log.Println(err)

	} else {
		// FIXME : Must be an implementation better than this.
		db.DB.Where("release_id = ?", c.Param("id")).Find(&rules)
		for i, _ := range rules {
			RulesID = rules[i].ID
			db.DB.Where("rule_id =?", RulesID).First(&buf)
			timerule = append(timerule, buf)

		}
		c.HTML(http.StatusOK, "release.html", gin.H{
			"id":       c.Param("id"),
			"release":  release,
			"timerule": timerule,
		})
	}

}

func EditRelease(c *gin.Context) {
	var release model.Release
	c.Bind(&release)

	db.DB.Model(&release).Where("id = ?", c.Param("id")).Updates(release)
	c.Redirect(http.StatusMovedPermanently, "/admin/dashboard/releases/")

}

func DelRelease(c *gin.Context) {
	var release model.Release

	db.DB.Where("id = ?", c.Param("id")).Delete(&release)
	c.Redirect(http.StatusMovedPermanently, "/admin/dashboard/releases/")
}

// New release
func NewRelease(c *gin.Context) {
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
		c.Redirect(http.StatusMovedPermanently, "/admin/dashboard")
	}
}

// Admin dashboard (new releases)
func Admin(c *gin.Context) {
	c.HTML(http.StatusOK, "dashboard.html", nil)
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
