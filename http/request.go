package http

import (
	"net/http"

	"code.videolan.org/GSoC2017/Marco/UpdateServer/core"
	"github.com/gin-gonic/gin"
)

func GetRequests(c *gin.Context) {
	requests := core.GetRequests(c.Param("channel"), c.Param("product"))
	c.HTML(http.StatusOK, "requests.html", gin.H{
		"requests": requests,
	})
}
