package main

import (
	"github.com/gin-gonic/gin"
	"github.com/xmarcoied/go-updater/model"
)

// In Process : Use Gin package
// Global variables are bad
var db model.Impl

func main() {
	// TODO : Add a middleware to keep the DB info
	db.ConnectDB()
	defer db.CloseDB()
	router := gin.Default()
	router.GET("/dashboard", admin)
	router.GET("/dashboard/showoff", adminshowoff)
	router.POST("/dashboard/new_release", newRelease)
	// TODO : status generation
	vlcRouter := router.Group("/vlc/:channel")
	{
		vlcRouter.GET("/status", updatesig)
		vlcRouter.GET("/showoff", showoff)
		vlcRouter.GET("/update", update)
	}
	router.Run(":80")
}

func ReleaseMap(r model.UpdateRequest) (model.Release, bool) {
	var ret model.Release
	var booleanRet bool
	ret = db.ReleaseMatch(r, ret)
	if booleanRet = true; ret.Channel == "" {
		booleanRet = false
	}
	return ret, booleanRet
}
