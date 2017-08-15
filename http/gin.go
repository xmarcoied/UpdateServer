package http

import (
	"encoding/json"
	"net/http"

	"code.videolan.org/GSoC2017/Marco/UpdateServer/core"
	"code.videolan.org/GSoC2017/Marco/UpdateServer/database"
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

func AddRelease(c *gin.Context) {
	channels := core.GetChannels()
	c.HTML(http.StatusOK, "newrelease.html", gin.H{
		"channels": channels,
	})
}

func NewRelease(c *gin.Context) {
	var (
		release database.Release
		buf     struct {
			Channel        string `json:"channel"`
			OS             string `json:"os"`
			OsVer          string `json:"os_ver"`
			OsArch         string `json:"os_arch"`
			ProductVersion string `json:"product_ver"`
			URL            string `json:"url"`
			Title          string `json:"title"`
			Description    string `json:"desc"`
			Product        string `json:"product"`
		}
	)
	c.Bind(&release)
	// FIXME: Isn't there a way to handle that?
	buf.Channel = release.Channel
	buf.OS = release.OS
	buf.OsVer = release.OsVer
	buf.OsArch = release.OsArch
	buf.ProductVersion = release.ProductVersion
	buf.URL = release.URL
	buf.Title = release.Title
	buf.Description = release.Description
	buf.Product = release.Product

	ReleaseJSON, _ := json.Marshal(buf)
	c.SetCookie("release", string(ReleaseJSON), 0, "/", "", false, false)
	c.Redirect(http.StatusMovedPermanently, "/admin/dashboard/addsignature/new")
}
