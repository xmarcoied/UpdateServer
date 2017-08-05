package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

	ReleaseJSON, _ := json.Marshal(release)
	c.SetCookie("release", string(ReleaseJSON), 0, "/", "", false, false)
	c.Redirect(http.StatusMovedPermanently, "/admin/dashboard/addsignature/edit?id="+c.Param("id"))
}

func (rlc ReleaseController) DelRelease(c *gin.Context) {
	var release model.Release

	db.DB.Where("id = ?", c.Param("id")).Delete(&release)
	c.Redirect(http.StatusMovedPermanently, "/admin/dashboard/releases/")
}

// New release
func (rlc ReleaseController) NewRelease(c *gin.Context) {

	var (
		release model.Release
		buf     struct {
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
	c.Bind(&release)
	// FIXME: Isn't there a way to handle that?
	buf.Channel = release.Channel
	buf.OS = release.OS
	buf.OsVer = release.OsVer
	buf.OsArch = release.OsArch
	buf.ProductVersion = release.ProductVersion
	buf.URL = release.URL
	buf.Title = release.Title
	buf.Description = release.Description
	buf.Product = release.Product

	ReleaseJSON, _ := json.Marshal(buf)
	c.SetCookie("release", string(ReleaseJSON), 0, "/", "", false, false)
	c.Redirect(http.StatusMovedPermanently, "/admin/dashboard/addsignature/new")
}

func (rlc ReleaseController) AddSignature(c *gin.Context) {
	release, _ := c.Cookie("release")
	var buf struct {
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

	json.Unmarshal([]byte(string(release)), &buf)

	ReleaseDir := "static/releases/tmp"
	ReleaseJSON, _ := json.Marshal(buf)
	ioutil.WriteFile(ReleaseDir, ReleaseJSON, 0644)

	ref := c.Param("reference")
	query := "0"
	if ref == "edit" {
		query = c.Query("id")
	}
	c.HTML(http.StatusOK, "newsignature.html", gin.H{
		"status": string(ReleaseJSON),
		"ref":    ref,
		"query":  query,
	})
}

func (rlc ReleaseController) VerifySignature(c *gin.Context) {
	var (
		binding struct {
			Content   string `form:"content" json:"content"`
			Signature string `form:"signature" json:"signature"`
		}
		release model.Release
	)

	c.Bind(&binding)
	json.Unmarshal([]byte(string(binding.Content)), &release)

	SignatureDir := "static/signatures/tmp.asc"
	ioutil.WriteFile(SignatureDir, []byte(binding.Signature), 0644)

	c.SetCookie("release", "", 0, "/", "", false, false)
	if utils.ProcessRelease(release) == true {
		if c.Param("reference") == "new" {
			db.DB.Create(&release)

		}
		if c.Param("reference") == "edit" {

			// This looks super ugly, I must improve it.
			id_buf, _ := strconv.Atoi(c.Query("id"))
			release.ID = uint(id_buf)
			db.DB.First(&release)
			var buf struct {
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
			json.Unmarshal([]byte(string(binding.Content)), &buf)
			release.Channel = buf.Channel
			release.OS = buf.OS
			release.OsVer = buf.OsVer
			release.OsArch = buf.OsArch
			release.ProductVersion = buf.ProductVersion
			release.URL = buf.URL
			release.Title = buf.Title
			release.Description = buf.Description
			release.Product = buf.Product

			db.DB.Save(&release)
		}

		os.Rename("static/releases/tmp", "static/releases/"+strconv.Itoa(int(release.ID)))
		os.Rename("static/signatures/tmp.asc", "static/signatures/"+strconv.Itoa(int(release.ID))+".asc")
		c.Redirect(http.StatusMovedPermanently, "/admin/dashboard/release/"+strconv.Itoa(int(release.ID)))

	} else {
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
