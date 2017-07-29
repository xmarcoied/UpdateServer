package controllers

import (
	"log"
	"net/http"

	"code.videolan.org/GSoC2017/Marco/UpdateServer/models"
	"code.videolan.org/GSoC2017/Marco/UpdateServer/utils"
	"github.com/gin-gonic/gin"
)

func NewRequestController() *RequestController {
	return &RequestController{}
}

// Show all requests
func (rc RequestController) GetRequests(c *gin.Context) {
	var requests []model.UpdateRequest
	requests = db.AllRequests(requests, c.Param("channel"), c.Param("product"))
	c.HTML(http.StatusOK, "requests.html", gin.H{
		"requests": requests,
	})
}

func (rc RequestController) Update(c *gin.Context) {
	// Request params are now getting in GET params
	var request model.UpdateRequest
	c.Bind(&request)

	request.Channel = c.Param("channel")
	request.Product = c.Param("product")

	ForwardHeader := c.Request.Header.Get("X-Forwarded-For")
	request.IP = utils.GetIP(ForwardHeader)

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
			if CheckOsRule(release) == false {
				return release, false
			}
			if CheckVersionRule(release) == false {
				return release, false
			}
			if CheckRollRule(release) == false {
				return release, false
			}
			found, check := CheckIPRule(release, r)
			if found == true {
				return release, check
			}

			return release, true
		}
		return ret, false
	}
}
