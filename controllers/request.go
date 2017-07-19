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
	var allAvailableReleases []model.Release
	var ret model.Release

	db.ReleaseMatch(r, &allAvailableReleases)

	// First check if there're any release matched with the request specs
	if len(allAvailableReleases) == 0 {
		return ret, false
	} else {

		// Check if the update request match the rules listed
		for _, release := range allAvailableReleases {
			if CountRules(release) == 0 {
				return release, true
			}
			if CheckTimeRule(release) == false {
				return release, false
			}

			return release, true
		}
		return ret, false
	}
}
