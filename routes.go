package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xmarcoied/go-updater/model"
	"log"
	"net/http"
)

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
