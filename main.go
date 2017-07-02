package main

import (
	"github.com/gin-gonic/gin"
	"github.com/xmarcoied/go-updater/model"
)

var db model.Impl

func main() {
	// TODO : Add a middleware to keep the DB info
	db.ConnectDB()
	defer db.CloseDB()
	router := gin.Default()

	adminRouter := router.Group("/admin")
	{
		adminRouter.GET("/dashboard", admin)
		adminRouter.GET("/dashboard/showoff", adminshowoff)
		adminRouter.POST("/dashboard/new_release", newRelease)
	}

	appRouter := router.Group("/u/:product/:channel")
	{
		appRouter.GET("/status", updatesig)
		appRouter.GET("/showoff", showoff)
		appRouter.GET("/update", update)
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
