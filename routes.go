package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xmarcoied/go-updater/model"
)

// Show all requests
func getRequests(c *gin.Context) {
	var requests []model.UpdateRequest
	requests = db.AllRequests(requests, c.Param("channel"), c.Param("product"))

	c.HTML(http.StatusOK, "requests.html", gin.H{
		"requests": requests,
	})
}

// Show all releases
func getReleases(c *gin.Context) {
	var releases []model.Release
	db.DB.Order("id").Find(&releases)

	c.HTML(http.StatusOK, "releases.html", gin.H{
		"releases": releases,
	})
}

func getRelease(c *gin.Context) {
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

func editRelease(c *gin.Context) {
	var release model.Release
	c.Bind(&release)

	db.DB.Model(&release).Where("id = ?", c.Param("id")).Updates(release)
	c.Redirect(http.StatusMovedPermanently, "/admin/dashboard/releases/")

}

func delRelease(c *gin.Context) {
	var release model.Release

	db.DB.Where("id = ?", c.Param("id")).Delete(&release)
	c.Redirect(http.StatusMovedPermanently, "/admin/dashboard/releases/")
}

// New release
func newRelease(c *gin.Context) {
	var release model.Release
	c.Bind(&release)
	log.Println(release)

	db.DB.Table("releases").Create(&release)
	c.Redirect(http.StatusMovedPermanently, "/admin/dashboard/releases/")
}

// Admin dashboard (new releases)
func admin(c *gin.Context) {
	c.HTML(http.StatusOK, "dashboard.html", nil)
}

func update(c *gin.Context) {
	// Request params are now getting in GET params
	var request model.UpdateRequest
	c.Bind(&request)
	request.Channel = c.Param("channel")
	request.Product = c.Param("product")

	matchedRelease, retStatus := ReleaseMap(request)
	if retStatus {
		log.Println("There's a release matched")
		request.Status = true
		c.JSON(http.StatusOK, matchedRelease)
	} else {
		log.Println("No release matched")
		request.Status = false
	}
	// TODO : DB Model API
	// FIXME : initiate the DB once and pass it everywhere
	db.DB.Table("update_requests").Create(&request)

}

func getChannels(c *gin.Context) {
	var channels []model.Channel
	db.DB.Order("id").Find(&channels)

	c.HTML(http.StatusOK, "channels.html", gin.H{
		"channels": channels,
	})

}

func addChannel(c *gin.Context) {
	c.HTML(http.StatusOK, "channel.html", nil)

}

func newChannel(c *gin.Context) {
	var channel model.Channel
	c.Bind(&channel)
	log.Println(channel)

	db.DB.Table("channels").Create(&channel)
	c.Redirect(http.StatusMovedPermanently, "/admin/dashboard/channels/")
}
