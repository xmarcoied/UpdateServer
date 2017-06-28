package main

import (
	"github.com/gin-gonic/gin"
	"github.com/xmarcoied/go-updater/model"
)

// In Process : Use Gin package
// Global variables are bad
var db model.Impl

func main() {
	db.ConnectDB()
	defer db.CloseDB()
	router := gin.Default()

	router.GET("/dashboard", admin)
	router.POST("/dashboard/new_release", newRelease)
	// TODO : status generation
	vlcRouter := router.Group("/vlc/:channel")
	{
		vlcRouter.StaticFile("status", "./client/static/status")
		vlcRouter.StaticFile("status.asc", "./client/static/status.asc")
		vlcRouter.GET("/showoff", showoff)
		vlcRouter.GET("/update", update)
	}
	router.Run(":80")
}
