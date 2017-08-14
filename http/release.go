package http

import (
	"net/http"

	"code.videolan.org/GSoC2017/Marco/UpdateServer/core"
	"github.com/gin-gonic/gin"
)

func GetReleases(c *gin.Context) {
	releases := core.GetReleases()
	c.HTML(http.StatusOK, "releases.html", gin.H{
		"releases": releases,
	})
}
