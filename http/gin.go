package http

import (
	"github.com/gin-gonic/gin"
)

func RouterInit() *gin.Engine {
	router := gin.Default()

	router.GET("/admin/dashboard/channels", GetChannels)
	router.GET("/admin/dashboard/channels/add", AddChannel)

	router.POST("/admin/dashboard/new_channel", NewChannel)

	router.LoadHTMLGlob("html/*.html")
	return router
}

func Run(addr string) {
	RouterInit().Run(":" + addr)
}
