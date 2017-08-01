package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"code.videolan.org/GSoC2017/Marco/UpdateServer/model"
	"code.videolan.org/GSoC2017/Marco/UpdateServer/utils"
	"github.com/gin-gonic/gin"
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
		release  model.Release
		rules    []model.Rule
		channels []model.Channel
	)
	if err := db.DB.Find(&release, "id = ?", c.Param("id")).Related(&rules).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		log.Println(err)

	} else {

		for i, rule := range rules {
			db.DB.Model(&rule).Related(&rules[i].TimeRule)
			db.DB.Model(&rule).Related(&rules[i].RollRule)
			db.DB.Model(&rule).Related(&rules[i].VersionRule)
			db.DB.Model(&rule).Related(&rules[i].OsRule)
			db.DB.Model(&rule).Related(&rules[i].IPRule)

		}

		db.DB.Model(&channels).Find(&channels)

		c.HTML(http.StatusOK, "release.html", gin.H{
			"release":  release,
			"rules":    rules,
			"channels": channels,
		})
	}

}

func (rlc ReleaseController) EditRelease(c *gin.Context) {
	var release model.Release
	c.Bind(&release)

	id_buf, _ := strconv.Atoi(c.Param("id"))
	release.ID = uint(id_buf)

	GenerateStatus(release)
	GenerateSignature(release)
	if utils.ProcessRelease(release) == true {
		log.Println("Accepted Release")
		db.DB.Model(&release).Where("id = ?", c.Param("id")).Updates(release)
		c.Redirect(http.StatusMovedPermanently, "/admin/dashboard/releases")

	} else {
		log.Println("Refused Release")
		db.DB.Delete(&release)
		c.Redirect(http.StatusMovedPermanently, "/admin/dashboard/releases")
	}
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

	// FIXME : if the connection dropped for any reason at this point
	// the server would count this release as a valid/signed release.

	db.DB.Save(&release)
	GenerateStatus(release)
	GenerateSignature(release)
	if utils.ProcessRelease(release) == true {
		log.Println("Accepted Release")
		c.Redirect(http.StatusMovedPermanently, "/admin/dashboard/releases/")

	} else {
		log.Println("Refused Release")
		db.DB.Delete(&release)
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
