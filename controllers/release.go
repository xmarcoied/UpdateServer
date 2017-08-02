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
	c.Redirect(http.StatusMovedPermanently, "/admin/dashboard/addsignature/"+strconv.Itoa(int(release.ID)))

}

func (rlc ReleaseController) AddSignature(c *gin.Context) {
	var (
		release model.Release
		buf     struct {
			ID             uint
			Channel        string `json:"channel"`
			OS             string `json:"os"`
			OsVer          string `json:"os_ver"`
			OsArch         string `json:"os_arch"`
			ProductVersion string `json:"product_ver"`
			URL            string `json:"url"`
			Title          string `json:"title"`
			Description    string `json:"desc"`
			Product        string `json:"product"`
		}
	)
	db.DB.Model(&release).First(&release, "id = ?", c.Param("id"))

	// FIXME: Isn't there a way to handle that?
	buf.ID = release.ID
	buf.Channel = release.Channel
	buf.OS = release.OS
	buf.OsVer = release.OsVer
	buf.OsArch = release.OsArch
	buf.ProductVersion = release.ProductVersion
	buf.URL = release.URL
	buf.Title = release.Title
	buf.Description = release.Description
	buf.Product = release.Product

	ReleaseDir := "static/releases/" + strconv.Itoa(int(release.ID))
	ReleaseJSON, _ := json.Marshal(buf)
	ioutil.WriteFile(ReleaseDir, ReleaseJSON, 0644)

	c.HTML(http.StatusOK, "newsignature.html", gin.H{
		"status": string(ReleaseJSON),
	})
}

func (rlc ReleaseController) VerifySignature(c *gin.Context) {
	var (
		buf struct {
			Content   string `form:"content" json:"content"`
			Signature string `form:"signature" json:"signature"`
		}
		release model.Release
	)

	c.Bind(&buf)
	json.Unmarshal([]byte(string(buf.Content)), &release)
	db.DB.Model(&release).First(&release, "id = ?", release.ID)

	SignatureDir := "static/signatures/" + strconv.Itoa(int(release.ID)) + ".asc"
	ioutil.WriteFile(SignatureDir, []byte(buf.Signature), 0644)

	if utils.ProcessRelease(release) == true {
		c.Redirect(http.StatusMovedPermanently, "/admin/dashboard/release/"+strconv.Itoa(int(release.ID)))

	} else {
		db.DB.Delete(&release)
		log.Println("Bad Signature")
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
