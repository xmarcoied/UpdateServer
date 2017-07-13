package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xmarcoied/go-updater/model"
)

// Show all requests
func GetRequests(c *gin.Context) {
	var requests []model.UpdateRequest
	requests = db.AllRequests(requests, c.Param("channel"), c.Param("product"))

	c.HTML(http.StatusOK, "requests.html", gin.H{
		"requests": requests,
	})
}

func Update(c *gin.Context) {
	// Request params are now getting in GET params
	var request model.UpdateRequest
	c.Bind(&request)
	request.Channel = c.Param("channel")
	request.Product = c.Param("product")

	matchedRelease, retStatus := ReleaseMap(request)
	if retStatus {
		log.Println("There's a release matched")
		request.Status = true
		c.JSON(http.StatusOK, matchedRelease)
	} else {
		log.Println("No release matched")
		request.Status = false
	}
	// TODO : DB Model API
	// FIXME : initiate the DB once and pass it everywhere
	db.DB.Table("update_requests").Create(&request)

}

func ReleaseMap(r model.UpdateRequest) (model.Release, bool) {
	var ret model.Release
	var booleanRet bool
	ret = db.ReleaseMatch(r, ret)
	if booleanRet = true; ret.Channel == "" {
		booleanRet = false
	}
	return ret, booleanRet
}
