package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq" // for database
	"time"
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
}

// Release database model
type Release struct {
	ID          uint      `gorm:"primary_key"`
	CreatedAt   time.Time `gorm:column:createdAt`
	Channel     string    `form:"channel"`
	OS          string    `form:"os"`
	OsVer       string    `form:"os_ver"`
	OsArch      string    `form:"os_arch"`
	VlcVer      string    `form:"vlc_ver"`
	URL         string    `form:"url"`
	Title       string    `form:"title"`
	Description string    `form:"desc"`
	Signature   string    `form:"sig"`
}

// Impl is handling gorm
type Impl struct {
	DB *gorm.DB
}

// ConnectDB initiate the database
func (i *Impl) ConnectDB() {
	// TODO : Move the psqlinfo to config & handle config/yml
	psqlInfo := "host=localhost dbname=marcoied user=postgres password=postgres sslmode=disable"
	i.DB, _ = gorm.Open("postgres", psqlInfo)
	i.DB.LogMode(true)
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

//AllRequests return all requests under specific channel
func (i *Impl) AllReleases(r []Release) []Release {
	i.DB.Find(&r)
	return r
}

//AllRequests return all requests under specific channel
func (i *Impl) AllRequests(r []UpdateRequest, ch string) []UpdateRequest {
	i.DB.Where("channel = ?", ch).Find(&r)
	return r
}

//CloseDB endup using the database
func (i *Impl) CloseDB() {
	i.DB.Close()
}

func (i *Impl) ReleaseMatch(req UpdateRequest, rel Release) Release {
	i.DB.Where("channel = ? AND os = ? AND os_arch = ? AND os_ver >= ?",
		req.Channel, req.OS, req.OsArch, req.OsVer).First(&rel)
	return rel
}
