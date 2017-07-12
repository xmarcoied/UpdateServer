package main

import (
	"flag"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/xmarcoied/go-updater/model"
)

var (
	db      model.Impl
	err     error
	addr    string
	ginMode string
	dbMode  string
)

func main() {
	// TODO : Add a middleware to keep the DB info

	flag.StringVar(&addr, "port", "8080", "a port to listen")
	flag.StringVar(&ginMode, "gin", "debug", "gin mode")
	flag.StringVar(&dbMode, "db", "false", "gin mode")

	flag.Parse()
	log.SetFlags(0)

	gin.SetMode(ginMode)
	err := db.ConnectDB(dbMode)
	if err != nil {
		log.Fatal(err)
	}

	RouterInit().Run(":" + addr)
}

func RouterInit() *gin.Engine {

	router := gin.Default()

	adminRouter := router.Group("/admin")
	{
		adminRouter.GET("/dashboard", admin)
		adminRouter.GET("/dashboard/releases", getReleases)
		adminRouter.GET("/dashboard/release/:id", getRelease)
		adminRouter.POST("/dashboard/new_release", newRelease)
		adminRouter.POST("/dashboard/edit_release/:id", editRelease)
	}

	appRouter := router.Group("/u/:product/:channel")
	{
		appRouter.GET("/get_requests", getRequests)
		appRouter.GET("/update", update)
	}
	router.LoadHTMLGlob("view/*.html")
	return router
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
