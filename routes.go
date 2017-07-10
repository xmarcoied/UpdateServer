package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xmarcoied/go-updater/model"
)

// Showoff all requests
func getRequests(c *gin.Context) {
	var requests []model.UpdateRequest
	requests = db.AllRequests(requests, c.Param("channel"), c.Param("product"))
	c.HTML(http.StatusOK, "requests.html", gin.H{
		"requests": requests,
	})
}

// Showoff all releases
func getReleases(c *gin.Context) {
	var releases []model.Release
	releases = db.AllReleases(releases)
	c.HTML(http.StatusOK, "releases.html", gin.H{
		"releases": releases,
	})
}

// New release
func newRelease(c *gin.Context) {
	var release model.Release
	c.Bind(&release)
	log.Println(release)

	db.NewRelease(release)
	c.Redirect(http.StatusMovedPermanently, "/admin/dashboard/get_releases/")
}

// Admin dashboard (new releases)
func admin(c *gin.Context) {
	t, _ := template.ParseFiles("view/dashboard.html")
	t.Execute(c.Writer, nil)
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
		c.JSON(200, matchedRelease)
	} else {
		log.Println("No release matched")
		request.Status = false
	}
	// TODO : DB Model API
	// FIXME : initiate the DB once and pass it everywhere
	db.NewRequest(request)

}
