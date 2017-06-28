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
	requests = db.AllRequests(requests, c.Param("channel"))

	for _, request := range requests {
		fmt.Fprint(c.Writer, "ID :", request.ID)
		fmt.Fprint(c.Writer, " , created_at :", request.CreatedAt)
		fmt.Fprint(c.Writer, " , OS :", request.OS)
		fmt.Fprint(c.Writer, " , OS_VER :", request.OsVer)
		fmt.Fprint(c.Writer, " , OS_ARCH :", request.OsArch)
		fmt.Fprint(c.Writer, " , VLC_VER :", request.VlcVer)
		fmt.Fprintln(c.Writer, " , IP :", request.IP)
	}

}

// New release
func newRelease(c *gin.Context) {
	var release model.Release
	c.Bind(&release)
	log.Println(release)

	db.NewRelease(release)
	redirectRoute := "http://update.videolan.org/dashboard"
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
	// TODO : DB Model API
	// FIXME : initiate the DB once and pass it everywhere
	db.NewRequest(request)
	statusRoute := "http://update.videolan.org/vlc/" + c.Param("channel") + "/status"
	c.Redirect(http.StatusMovedPermanently, statusRoute)

}
