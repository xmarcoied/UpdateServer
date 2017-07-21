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
	if err != nil {
		log.Fatal(err)
	}

	RouterInit().Run(":" + addr)
}

func RouterInit() *gin.Engine {

	router := gin.Default()

	adminRouter := router.Group("/admin")
	{
		adminRouter.GET("/dashboard/newrelease", controllers.AddRelease)
		adminRouter.GET("/dashboard/releases", controllers.GetReleases)
		adminRouter.GET("/dashboard/release/:id", controllers.GetRelease)
		adminRouter.GET("/dashboard/del_release/:id", controllers.DelRelease)

		adminRouter.POST("/dashboard/new_release", controllers.NewRelease)
		adminRouter.POST("/dashboard/edit_release/:id", controllers.EditRelease)

		adminRouter.GET("/dashboard/channels", controllers.GetChannels)
		adminRouter.GET("/dashboard/channels/add", controllers.AddChannel)

		adminRouter.POST("/dashboard/new_channel", controllers.NewChannel)

		adminRouter.GET("/dashboard/add_rule/:id", controllers.AddRule)
		adminRouter.POST("/dashboard/new_rule/:rule/:id", controllers.NewRule)

	}

	appRouter := router.Group("/u/:product/:channel")
	{
		appRouter.GET("/get_requests", controllers.GetRequests)
		appRouter.GET("/update", controllers.Update)
	}

	router.LoadHTMLGlob("view/*.html")
	return router
}
