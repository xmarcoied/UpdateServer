package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xmarcoied/go-updater/model"
	"html/template"
	"log"
	"net/http"
)

// Showoff all requests
func showoff(c *gin.Context) {
	var requests []model.UpdateRequest
	requests = db.AllRequests(requests, c.Param("channel"), c.Param("product"))
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
	redirectRoute := "http://update.videolan.org/admin/dashboard/showoff"
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
	request.Channel = c.Param("channel")
	request.Product = c.Param("product")

	matchedRelease, retStatus := ReleaseMap(request)
	if retStatus {
		log.Println("There's a release matched")
		request.Status = true
		// First Line : version number
		fmt.Fprintln(c.Writer, matchedRelease.VlcVer)
		// Second Line : URL
		fmt.Fprintln(c.Writer, matchedRelease.URL)
		// Third Line : Desc.
		fmt.Fprintln(c.Writer, matchedRelease.Title)
		fmt.Fprint(c.Writer, matchedRelease.Description)

	} else {
		log.Println("No release matched")
		request.Status = false
	}
	// TODO : DB Model API
	// FIXME : initiate the DB once and pass it everywhere
	db.NewRequest(request)

}

func updatesig(c *gin.Context) {
	var request model.UpdateRequest
	c.Bind(&request)
	request.Channel = c.Param("channel")
	request.Product = c.Param("product")

	matchedRelease, retStatus := ReleaseMap(request)
	if retStatus {
		log.Println("There's a release matched")
		fmt.Fprintln(c.Writer, matchedRelease.Signature)

	} else {
		log.Println("No release matched")
	}
}
