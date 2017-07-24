package main

import (
	"flag"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/xmarcoied/go-updater/controllers"
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
	controllers.SetDB(db)
	rc := controllers.NewRequestController()
	cc := controllers.NewChannelController()
	rsc := controllers.NewRulesController()
	rlc := controllers.NewReleaseController()

	if err != nil {
		log.Fatal(err)
	}

	RouterInit(rc, cc, rsc, rlc).Run(":" + addr)
}

func RouterInit(rc *controllers.RequestController, cc *controllers.ChannelController, rsc *controllers.RulesController, rlc *controllers.ReleaseController) *gin.Engine {

	router := gin.Default()

	adminRouter := router.Group("/admin")
	{
		adminRouter.GET("/dashboard/newrelease", rlc.AddRelease)
		adminRouter.GET("/dashboard/releases", rlc.GetReleases)
		adminRouter.GET("/dashboard/release/:id", rlc.GetRelease)
		adminRouter.GET("/dashboard/del_release/:id", rlc.DelRelease)

		adminRouter.POST("/dashboard/new_release", rlc.NewRelease)
		adminRouter.POST("/dashboard/edit_release/:id", rlc.EditRelease)

		adminRouter.GET("/dashboard/channels", cc.GetChannels)
		adminRouter.GET("/dashboard/channels/add", cc.AddChannel)

		adminRouter.POST("/dashboard/new_channel", cc.NewChannel)

		adminRouter.GET("/dashboard/add_rule/:id", rsc.AddRule)
		adminRouter.POST("/dashboard/new_rule/:rule/:id", rsc.NewRule)

	}

	appRouter := router.Group("/u/:product/:channel")
	{
		appRouter.GET("/get_requests", rc.GetRequests)
		appRouter.GET("/update", rc.Update)
	}

	router.LoadHTMLGlob("view/*.html")
	return router
}
