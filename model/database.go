package model

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq" // for database
)

// UpdateRequest database model
type UpdateRequest struct {
	//gorm.Model
	ID        uint      `gorm:"primary_key"`
	CreatedAt time.Time `gorm:column:createdAt`
	Channel   string    `form:"channel"`
	OS        string    `form:"os"`
	OsVer     string    `form:"os_ver"`
	OsArch    string    `form:"os_arch"`
	VlcVer    string    `form:"vlc_ver"`
	IP        string    `form:"ip"`
	Status    bool
	Product   string
}

// Release database model
type Release struct {
	ID          uint      `gorm:"primary_key"`
	CreatedAt   time.Time `gorm:column:createdAt`
	Channel     string    `form:"channel" json:"channel"`
	OS          string    `form:"os" json:"os"`
	OsVer       string    `form:"os_ver" json:"os_ver"`
	OsArch      string    `form:"os_arch" json:"os_arch"`
	VlcVer      string    `form:"vlc_ver" json:"vlc_ver"`
	URL         string    `form:"url" json:"url"`
	Title       string    `form:"title" json:"title"`
	Description string    `form:"desc" json:"desc`
	Signature   string    `form:"sig" json:"sig"`
	Product     string    `form:"product" json:"product"`
}

// Impl is handling gorm
type Impl struct {
	DB *gorm.DB
}

// ConnectDB initiate the database
func (i *Impl) ConnectDB(ginMode string) {
	// TODO : Move the psqlinfo to config & handle config/yml
	psqlInfo := "host=localhost dbname=marcoied user=postgres password=postgres sslmode=disable"
	i.DB, _ = gorm.Open("postgres", psqlInfo)
	if ginMode == "true" {
		i.DB.LogMode(true)
	}
	i.DB.AutoMigrate(&UpdateRequest{}, &Release{})
}

// NewRelease add/create new update release
func (i *Impl) NewRelease(r Release) {
	i.DB.Create(&r)
}

// NewRequest add/create new update request
func (i *Impl) NewRequest(r UpdateRequest) {
	i.DB.Create(&r)
}

//AllReleases return all requests under specific channel
func (i *Impl) AllReleases(r []Release) []Release {
	i.DB.Find(&r)
	return r
}

//AllRequests return all requests under specific channel
func (i *Impl) AllRequests(r []UpdateRequest, ch string, p string) []UpdateRequest {
	i.DB.Where("channel = ? AND product = ?", ch, p).Find(&r)
	return r
}

//CloseDB endup using the database
func (i *Impl) CloseDB() {
	i.DB.Close()
}

func (i *Impl) ReleaseMatch(req UpdateRequest, rel Release) Release {
	i.DB.Where("product = ? AND channel = ? AND os = ? AND os_arch = ? AND os_ver >= ?",
		req.Product, req.Channel, req.OS, req.OsArch, req.OsVer).First(&rel)
	return rel
}
