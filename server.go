package main

import (
	"encoding/base64"
	"flag"
	"log"
	"math/rand"
	"net/http"
	"strings"
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

	auth := router.Group("/", Auth)
	adminRouter := auth.Group("/admin")
	{
		adminRouter.GET("/dashboard/newrelease", rlc.AddRelease)
		adminRouter.GET("/dashboard/releases", rlc.GetReleases)
		adminRouter.GET("/dashboard/release/:id", rlc.GetRelease)
		adminRouter.GET("/dashboard/del_release/:id", rlc.DelRelease)
		adminRouter.GET("/dashboard/addsignature/:reference", rlc.AddSignature)

		adminRouter.POST("/dashboard/verifysignature/:reference", rlc.VerifySignature)
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
		appRouter.GET("/requests", rc.GetRequests)
		appRouter.GET("/update", rc.Update)
		appRouter.GET("/signature", rc.GetSignature)
	}

	router.LoadHTMLGlob("view/*.html")
	return router
}

func Auth(c *gin.Context) {
	if checkAuth(c) {
		c.Next()
	} else {
		c.Writer.Header().Set("WWW-Authenticate", "Basic realm=UpdateServer")
		c.AbortWithStatus(http.StatusUnauthorized)

	}

}
func checkAuth(c *gin.Context) bool {
	auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)
	if len(auth) != 2 {
		return false
	}

	base, err := base64.StdEncoding.DecodeString(auth[1])
	if err != nil {
		return false
	}

	UserData := strings.SplitN(string(base), ":", 2)
	if len(UserData) != 2 {
		return false
	}
	log.Println(UserData)
	username := UserData[0]
	password := UserData[1]
	return username == "admin" && password == "admin"
}
