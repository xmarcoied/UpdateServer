package model

import (
	"time"

	"code.videolan.org/GSoC2017/Marco/UpdateServer/config"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq" // for database
)

// UpdateRequest database model
type UpdateRequest struct {
	//gorm.Model
	CreatedAt      time.Time
	ID             uint   `gorm:"primary_key"`
	Channel        string `form:"channel"`
	OS             string `form:"os"`
	OsVer          string `form:"os_ver"`
	OsArch         string `form:"os_arch"`
	ProductVersion string `form:"product_ver"`
	IP             string `form:"ip"`
	Status         bool
	Product        string
}

// Release database model
type Release struct {
	CreatedAt      time.Time
	UpdatedAt      time.Time
	ID             uint   `gorm:"primary_key"`
	Channel        string `form:"channel" json:"channel"`
	OS             string `form:"os" json:"os"`
	OsVer          string `form:"os_ver" json:"os_ver"`
	OsArch         string `form:"os_arch" json:"os_arch"`
	ProductVersion string `form:"product_ver" json:"product_ver"`
	URL            string `form:"url" json:"url"`
	Title          string `form:"title" json:"title"`
	Description    string `form:"desc" json:"desc"`
	Product        string `form:"product" json:"product"`
	Rules          Rule
	Signature      string
}

// Channel database model
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
func (i *Impl) ConnectDB(c *config.Configuration) error {
	var err error
	// TODO : Move the psqlinfo to config & handle config/yml
	psqlInfo := "host=" + c.Database.Host + " dbname=" + c.Database.Name + " user=" + c.Database.User + " password=" + c.Database.Password + " sslmode=disable"
	i.DB, err = gorm.Open("postgres", psqlInfo)
	i.DB.LogMode(true)
	i.DB.AutoMigrate(&UpdateRequest{}, &Release{}, &Channel{}, &Rule{}, &TimeRule{}, &OsRule{}, &VersionRule{}, &IPRule{}, &RollRule{})

	return err
}

//AllRequests return all requests under specific channel
func (i *Impl) AllRequests(r []UpdateRequest, ch string, p string) []UpdateRequest {
	i.DB.Table("update_requests").Where("product = ? AND channel = ?", p, ch).Find(&r)

	return r
}

func (i *Impl) ReleaseMatch(req UpdateRequest, rel *[]Release) {
	i.DB.Where("product = ? AND channel = ? AND os = ? AND os_arch = ? AND os_ver >= ?",
		req.Product, req.Channel, req.OS, req.OsArch, req.OsVer).Find(&rel)

}
