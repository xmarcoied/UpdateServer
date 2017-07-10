package main

import (
	"flag"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/xmarcoied/go-updater/model"
)

var (
	db   model.Impl
	addr string
)

func main() {
	// TODO : Add a middleware to keep the DB info

	flag.StringVar(&addr, "port", "8080", "a port to listen")
	flag.Parse()
	log.SetFlags(0)
	db.ConnectDB()
	defer db.CloseDB()
	router := gin.Default()

	adminRouter := router.Group("/admin")
	{
		adminRouter.GET("/dashboard", admin)
		adminRouter.GET("/dashboard/get_releases", getReleases)
		adminRouter.POST("/dashboard/new_release", newRelease)
	}

	appRouter := router.Group("/u/:product/:channel")
	{
		appRouter.GET("/get_requests", getRequests)
		appRouter.GET("/update", update)
	}

	router.LoadHTMLGlob("view/*.html")
	router.Run(":" + addr)
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
