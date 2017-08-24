package http

import (
	"fmt"
	"net/http"
	"strings"

	"code.videolan.org/GSoC2017/Marco/UpdateServer/core"
	"code.videolan.org/GSoC2017/Marco/UpdateServer/database"
	"github.com/gin-gonic/gin"
)

// GetRequests is http handler to represent all the requests available in the UpdateServer
func GetRequests(c *gin.Context) {
	var (
		query   string
		request struct {
			Channel string `form:"channel"`
			Product string `form:"product"`
		}
	)
	c.Bind(&request)
	if request.Channel != "" {
		newquery := fmt.Sprintf("channel = '%s'", request.Channel)
		query = database.QueryAppend(query, newquery)
	}

	if request.Product != "" {
		newquery := fmt.Sprintf("product = '%s'", request.Product)
		query = database.QueryAppend(query, newquery)
	}

	requests := core.GetRequests(query)
	channels := core.GetChannels()
	c.HTML(http.StatusOK, "requests.html", gin.H{
		"requests": requests,
		"channels": channels,
		"request":  request,
	})
}

// GetSignature
// Expecting requests like "x.asc" format
func GetSignature(c *gin.Context) {
	var release_id int
	fmt.Sscanf(c.Query("id"), "%d.asc", &release_id)
	releaseSignature := core.GetSignature(release_id)
	c.String(http.StatusOK, releaseSignature)

}

func Update(c *gin.Context) {
	// Request params are now getting in GET params
	var request database.UpdateRequest
	c.Bind(&request)

	ForwardHeader := c.Request.Header.Get("X-Forwarded-For")
	ForwardedFields := strings.Split(ForwardHeader, ",")
	if len(ForwardedFields) > 0 {
		request.IP = ForwardedFields[0]
	} else {
		request.IP = ""
	}

	matchedRelease, retStatus := core.ReleaseMap(request)
	if retStatus {
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

func GetAct(c *gin.Context) {
	var request database.UpdateRequest
	c.Bind(&request)

	channels := core.GetChannels()
	release, _ := core.ReleaseMap(request)
	c.HTML(http.StatusOK, "requestclient.html", gin.H{
		"request":  request,
		"release":  release,
		"channels": channels,
	})
}
