package main

import (
	"flag"
	"log"
	"math/rand"
	"time"

	"code.videolan.org/GSoC2017/Marco/UpdateServer/config"
	"code.videolan.org/GSoC2017/Marco/UpdateServer/controllers"
	"code.videolan.org/GSoC2017/Marco/UpdateServer/model"
	"github.com/gin-gonic/gin"
)

var (
	db   model.Impl
	addr string
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	flag.StringVar(&addr, "port", "8080", "The port server will be running on")
	flag.Parse()
}

func main() {

	c, err := config.Get()
	if err != nil {
		log.Fatal(err)
	}

	err = db.ConnectDB(c)
	if err != nil {
		log.Fatal(err)
	}

	controllers.SetDB(db)
	rc := controllers.NewRequestController()
	cc := controllers.NewChannelController()
	rsc := controllers.NewRulesController()
	rlc := controllers.NewReleaseController()

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
		adminRouter.GET("/dashboard/delete_rule/:rule/:id", rsc.DeleteRule)

	}

	appRouter := router.Group("/u/:product/:channel")
	{
		appRouter.GET("/get_requests", rc.GetRequests)
		appRouter.GET("/update", rc.Update)
	}

	router.LoadHTMLGlob("view/*.html")
	return router
}
