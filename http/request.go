package http

import (
	"log"
	"net/http"

	"code.videolan.org/GSoC2017/Marco/UpdateServer/core"
	"code.videolan.org/GSoC2017/Marco/UpdateServer/database"
	"github.com/gin-gonic/gin"
)

// GetRequests is http handler to represent all the requests available in the UpdateServer
func GetRequests(c *gin.Context) {
	requests := core.GetRequests(c.Param("channel"), c.Param("product"))
	c.HTML(http.StatusOK, "requests.html", gin.H{
		"requests": requests,
	})
}

// GetSignature
func GetSignature(c *gin.Context) {
	releaseSignature := core.GetSignature(c.Query("id"))
	c.String(http.StatusOK, releaseSignature)
}

func Update(c *gin.Context) {
	// Request params are now getting in GET params
	var request database.UpdateRequest
	c.Bind(&request)

	request.Channel = c.Param("channel")
	request.Product = c.Param("product")

	matchedRelease, retStatus := ReleaseMap(request)
	if retStatus {
		log.Println("There's a release matched")
		request.Status = true
		var buf struct {
			ID             uint   `json:"id"`
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
		buf.ID = matchedRelease.ID
		buf.Channel = matchedRelease.Channel
		buf.OS = matchedRelease.OS
		buf.OsVer = matchedRelease.OsVer
		buf.OsArch = matchedRelease.OsArch
		buf.ProductVersion = matchedRelease.ProductVersion
		buf.URL = matchedRelease.URL
		buf.Title = matchedRelease.Title
		buf.Description = matchedRelease.Description
		buf.Product = matchedRelease.Product

		c.JSON(http.StatusOK, buf)
	} else {
		log.Println("No release matched")
		var buf struct {
			ID             uint   `json:"id"`
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
		buf.ID = matchedRelease.ID
		buf.Channel = matchedRelease.Channel
		buf.OS = matchedRelease.OS
		buf.OsVer = matchedRelease.OsVer
		buf.OsArch = matchedRelease.OsArch
		buf.ProductVersion = matchedRelease.ProductVersion
		buf.URL = matchedRelease.URL
		buf.Title = matchedRelease.Title
		buf.Description = matchedRelease.Description
		buf.Product = matchedRelease.Product

		c.JSON(http.StatusOK, buf)
		request.Status = false
	}

	core.NewRequest(request)
}

func ReleaseMap(request database.UpdateRequest) (database.Release, bool) {
	var release database.Release
	return release, true
}
