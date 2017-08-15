package http

import (
	"github.com/gin-gonic/gin"
)

// RouterInit function initiate the gin engine main handler
func RouterInit() *gin.Engine {
	router := gin.Default()

	router.GET("/admin/dashboard/channels", GetChannels)
	router.GET("/admin/dashboard/releases", GetReleases)
	router.GET("/admin/dashboard/release/:id", GetRelease)
	router.GET("/admin/dashboard/channels/add", AddChannel)
	router.GET("/u/:product/:channel/requests", GetRequests)
	router.GET("/admin/dashboard/newrelease", AddRelease)
	router.GET("/admin/dashboard/addsignature/:reference", AddSignature)

	router.POST("/admin/dashboard/edit_release/:id", EditRelease)
	router.POST("/admin/dashboard/verifysignature/:reference", VerifySignature)
	router.POST("/admin/dashboard/new_release", NewRelease)
	router.POST("/admin/dashboard/new_channel", NewChannel)

	router.LoadHTMLGlob("html/*.html")
	return router
}

// Run function initaite the UpdateServer
func Run(addr string) {
	RouterInit().Run(":" + addr)
}
