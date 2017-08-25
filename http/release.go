package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"code.videolan.org/GSoC2017/Marco/UpdateServer/core"
	"code.videolan.org/GSoC2017/Marco/UpdateServer/database"
	"code.videolan.org/GSoC2017/Marco/UpdateServer/utils"
	"github.com/gin-gonic/gin"
)

// GetReleases is http handler to represent all the releases available in the UpdateServer
func GetReleases(c *gin.Context) {
	// pass empty query to get all releases
	query := ""
	releases := core.GetReleases(query)
	c.HTML(http.StatusOK, "releases.html", gin.H{
		"releases": releases,
	})
}

// GetRelease
func GetRelease(c *gin.Context) {
	release := core.GetRelease(c.Param("id"))
	channels := core.GetChannels()
	rules := core.GetRules(release)
	c.HTML(http.StatusOK, "release.html", gin.H{
		"release":  release,
		"rules":    rules,
		"channels": channels,
	})
}

// ToggleRelease
func ToggleRelease(c *gin.Context) {
	release := core.GetRelease(c.Param("id"))
	core.ToggleReleaseActivtion(&release)
	c.Redirect(http.StatusMovedPermanently, "/admin/dashboard/release/"+c.Param("id"))
}

// AddRelease
func AddRelease(c *gin.Context) {
	channels := core.GetChannels()
	c.HTML(http.StatusOK, "newrelease.html", gin.H{
		"channels": channels,
	})
}

// NewRelease
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
	ReleaseChannel := core.GetChannel(release.Channel)
	c.SetCookie("release", string(ReleaseJSON), 0, "/", "", false, false)

	// Check if the channel have private key or not
	if ReleaseChannel.PrivateKey == "" {
		c.Redirect(http.StatusMovedPermanently, "/admin/dashboard/addsignature/new")

	} else {
		signature, err := utils.Sign(ReleaseChannel, release, string(ReleaseJSON))
		if err != nil {
			log.Println(err)
		}

		isvalid, err := utils.ProcessRelease(ReleaseChannel, release, signature, string(ReleaseJSON))
		log.Println(isvalid, err)
		if isvalid == true && err == nil {
			release.Signature = signature
			core.NewRelease(&release)
			c.Redirect(http.StatusMovedPermanently, "/admin/dashboard/release/"+strconv.Itoa(int(release.ID)))

		} else {
			log.Println(err)
			c.Redirect(http.StatusMovedPermanently, "/admin/dashboard/newrelease")
		}
	}
}

//DelRelease
func DelRelease(c *gin.Context) {
	core.DeleteRelease(c.Param("id"))
}

// EditRelease
func EditRelease(c *gin.Context) {
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
	ReleaseChannel := core.GetChannel(release.Channel)
	c.SetCookie("release", string(ReleaseJSON), 0, "/", "", false, false)

	// Check if the channel have private key or not
	if ReleaseChannel.PrivateKey == "" {
		c.Redirect(http.StatusMovedPermanently, "/admin/dashboard/addsignature/edit?id="+c.Param("id"))
	} else {
		signature, err := utils.Sign(ReleaseChannel, release, string(ReleaseJSON))
		if err != nil {
			log.Println(err)
		}

		isvalid, err := utils.ProcessRelease(ReleaseChannel, release, signature, string(ReleaseJSON))
		if isvalid == true && err == nil {
			core.EditRelease(&release, c.Param("id"), signature, string(ReleaseJSON))
			c.Redirect(http.StatusMovedPermanently, "/admin/dashboard/release/"+c.Param("id"))

		} else {
			log.Println(err)
			c.Redirect(http.StatusMovedPermanently, "/admin/dashboard/newrelease")
		}
	}
}

//AddSignature
func AddSignature(c *gin.Context) {
	release, _ := c.Cookie("release")
	var buf struct {
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

	json.Unmarshal([]byte(string(release)), &buf)

	ReleaseJSON, _ := json.Marshal(buf)

	ref := c.Param("reference")
	query := "0"
	if ref == "edit" {
		query = c.Query("id")
	}

	channel := core.GetChannel(buf.Channel)
	fingerprint, _ := utils.GetFingerprint(channel)
	status := fmt.Sprintf("echo -n '%s' | gpg --default-key %s --detach-sign -a", string(ReleaseJSON), fingerprint)
	c.HTML(http.StatusOK, "newsignature.html", gin.H{
		"status": status,
		"ref":    ref,
		"query":  query,
	})
}

// VerifySignature
func VerifySignature(c *gin.Context) {
	var (
		binding struct {
			Content   string `form:"content" json:"content"`
			Signature string `form:"signature" json:"signature"`
		}
		release database.Release
	)

	c.Bind(&binding)
	binding.Content, _ = c.Cookie("release")
	json.Unmarshal([]byte(string(binding.Content)), &release)
	channel := core.GetChannel(release.Channel)

	c.SetCookie("release", "", 0, "/", "", false, false)
	isvalid, err := utils.ProcessRelease(channel, release, binding.Signature, binding.Content)
	if isvalid == true && err == nil {
		if c.Param("reference") == "new" {
			release.Signature = binding.Signature
			core.NewRelease(&release)
		}
		if c.Param("reference") == "edit" {
			core.EditRelease(&release, c.Query("id"), binding.Signature, binding.Content)
		}
		c.Redirect(http.StatusMovedPermanently, "/admin/dashboard/release/"+strconv.Itoa(int(release.ID)))

	} else {
		c.Redirect(http.StatusMovedPermanently, "/admin/dashboard/newrelease")
	}
}
