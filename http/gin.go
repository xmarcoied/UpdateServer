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
	router.GET("/u/:product/:channel/signature", GetSignature)
	// router.GET("/u/:product/:channel/update", Update)
	router.GET("/admin/dashboard/newrelease", AddRelease)
	router.GET("/admin/dashboard/addsignature/:reference", AddSignature)
	router.GET("/admin/dashboard/del_release/:id", DelRelease)
	router.GET("/admin/dashboard/add_rule/:id", AddRule)
	router.GET("/admin/dashboard/delete_rule/:rule/:id", DeleteRule)

	router.POST("/admin/dashboard/new_rule/:rule/:id", NewRule)
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
