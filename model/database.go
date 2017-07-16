package model

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq" // for database
)

// UpdateRequest database model
type UpdateRequest struct {
	//gorm.Model
	CreatedAt time.Time
	ID        uint   `gorm:"primary_key"`
	Channel   string `form:"channel"`
	OS        string `form:"os"`
	OsVer     string `form:"os_ver"`
	OsArch    string `form:"os_arch"`
	VlcVer    string `form:"vlc_ver"`
	IP        string `form:"ip"`
	Status    bool
	Product   string
}

// Release database model
type Release struct {
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ID          uint   `gorm:"primary_key"`
	Channel     string `form:"channel" json:"channel"`
	OS          string `form:"os" json:"os"`
	OsVer       string `form:"os_ver" json:"os_ver"`
	OsArch      string `form:"os_arch" json:"os_arch"`
	VlcVer      string `form:"vlc_ver" json:"vlc_ver"`
	URL         string `form:"url" json:"url"`
	Title       string `form:"title" json:"title"`
	Description string `form:"desc" json:"desc"`
	Signature   string
	Product     string `form:"product" json:"product"`
}

type Channel struct {
	ID            uint   `gorm:"primary_key"`
	Name          string `form:"name" json:"name"`
	PublicKey     string `form:"pubkey" json:"pubkey"`
	PrivateKey    string `form:"privatekey" json:"privatekey"`
	ReleasesCount string
	RequestsCount string
}

// Impl is handling gorm
type Impl struct {
	DB *gorm.DB
}

// ConnectDB initiate the database
func (i *Impl) ConnectDB(dbMode string) error {
	var err error
	// TODO : Move the psqlinfo to config & handle config/yml
	psqlInfo := "host=localhost dbname=marcoied user=postgres password=postgres sslmode=disable"
	i.DB, err = gorm.Open("postgres", psqlInfo)
	if dbMode == "true" {
		i.DB.LogMode(true)
	}

	i.DB.AutoMigrate(&UpdateRequest{}, &Release{}, &Channel{})

	return err
}

//AllRequests return all requests under specific channel
func (i *Impl) AllRequests(r []UpdateRequest, ch string, p string) []UpdateRequest {
	i.DB.Table("update_requests").Where("product = ? AND channel = ?", p, ch).Find(&r)

	return r
}

func (i *Impl) ReleaseMatch(req UpdateRequest, rel Release) Release {
	i.DB.Table("releases").Where("product = ? AND channel = ? AND os = ? AND os_arch = ? AND os_ver >= ?",
		req.Product, req.Channel, req.OS, req.OsArch, req.OsVer).First(&rel)

	return rel
}
