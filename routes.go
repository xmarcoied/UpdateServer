package main

import (
	"github.com/gin-gonic/gin"
	"github.com/xmarcoied/go-updater/model"
	"html/template"
	"log"
	"net/http"
)

// Showoff all requests
func showoff(c *gin.Context) {
	var requests []model.UpdateRequest
	requests = db.AllRequests(requests, c.Param("channel"))
	t, _ := template.ParseFiles("view/requests.html")
	t.Execute(c.Writer, requests)

}

// Showoff all releases
func adminshowoff(c *gin.Context) {
	var requests []model.Release
	requests = db.AllReleases(requests)
	t, _ := template.ParseFiles("view/releases.html")
	t.Execute(c.Writer, requests)
}

// New release
func newRelease(c *gin.Context) {
	var release model.Release
	c.Bind(&release)
	log.Println(release)

	db.NewRelease(release)
	redirectRoute := "http://update.videolan.org/dashboard/showoff"
	c.Redirect(http.StatusMovedPermanently, redirectRoute)
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
	log.Println(request)
	request.Channel = c.Param("channel")

	_, retStatus := ReleaseMap(request)
	if retStatus {
		log.Println("There's a release matched")
		request.Status = true
	} else {
		log.Println("No release matched")
		request.Status = false
	}
	// TODO : DB Model API
	// FIXME : initiate the DB once and pass it everywhere
	db.NewRequest(request)
	statusRoute := "http://update.videolan.org/vlc/" + c.Param("channel") + "/status"
	c.Redirect(http.StatusMovedPermanently, statusRoute)

}
