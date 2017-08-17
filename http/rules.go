package http

import (
	"net/http"
	"strconv"
	"time"

	"code.videolan.org/GSoC2017/Marco/UpdateServer/core"
	"code.videolan.org/GSoC2017/Marco/UpdateServer/database"
	"github.com/gin-gonic/gin"
)

func AddRule(c *gin.Context) {
	c.HTML(http.StatusOK, "newrule.html", gin.H{
		"id": c.Param("id"),
	})
}

func NewRule(c *gin.Context) {

	var rule database.Rule
	switch c.Param("rule") {
	case "time":
		var buf struct {
			Start string `form:"timestart"`
			End   string `form:"timeend"`
		}
		c.Bind(&buf)

		layout := "2006-01-02T15:04"
		t_start, _ := time.Parse(layout, buf.Start)
		t_end, _ := time.Parse(layout, buf.End)

		if t_end.IsZero() {
			t_end, _ = time.Parse(layout, "2906-01-02T15:04")
		}
		rule.TimeRule.StartTime = t_start
		rule.TimeRule.EndTime = t_end

	case "os":
		var buf struct {
			OsVersion string `form:"os_version"`
		}
		c.Bind(&buf)
		rule.OsRule.OsVersion = buf.OsVersion

	case "version":
		var buf struct {
			ProductVersion string `form:"product_version"`
		}
		c.Bind(&buf)
		rule.VersionRule.ProductVersion = buf.ProductVersion

	case "ip":
		var buf struct {
			IP string `form:"ip_address"`
		}
		c.Bind(&buf)
		rule.IPRule.IP = buf.IP

	case "roll":
		var buf struct {
			RollingPercentage int `form:"roll"`
		}
		c.Bind(&buf)
		rule.RollRule.RollingPercentage = buf.RollingPercentage
	}

	core.NewRule(c.Param("id"), rule)
	c.Redirect(http.StatusMovedPermanently, "/admin/dashboard/release/"+c.Param("id"))
}

func DeleteRule(c *gin.Context) {
	release_id := strconv.Itoa(core.DeleteRule(c.Param("id")))
	c.Redirect(http.StatusMovedPermanently, "/admin/dashboard/release/"+release_id)

}
