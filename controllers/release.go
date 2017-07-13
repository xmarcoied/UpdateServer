package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xmarcoied/go-updater/model"
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
	var release model.Release

	if err := db.DB.Where("id = ?", c.Param("id")).First(&release).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		log.Println(err)

	} else {
		c.HTML(http.StatusOK, "release.html", gin.H{
			"release": release,
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
	log.Println(release)

	db.DB.Table("releases").Create(&release)
	c.Redirect(http.StatusMovedPermanently, "/admin/dashboard/releases/")
}

// Admin dashboard (new releases)
func Admin(c *gin.Context) {
	c.HTML(http.StatusOK, "dashboard.html", nil)
}
