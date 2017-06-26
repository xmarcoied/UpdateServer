package model

import(
//	"log"
	"time"
	"github.com/jinzhu/gorm"
       _"github.com/lib/pq"
)

type Update_Request struct{
	//gorm.Model
	ID          uint       `gorm:"primary_key"`
	Created_At  time.Time  `gorm:column:createdAt`
	Channel		string 	   `form:"channel"`
	OS 	 		string	   `form:"os"` 	   
	OS_VER 		string	   `form:"os_ver"`
	OS_ARCH 	string 	   `form:"os_arch"`
	VLC_VER 	string	   `form:"vlc_ver"`
	IP 			string     `form:"ip"`
}

type Impl struct {
    	DB *gorm.DB
}

func (i *Impl) ConnectDB(){
	// TODO : Move the psqlinfo to config & handle config/yml
	psqlInfo := "host=localhost dbname=marcoied user=postgres password=postgres sslmode=disable"
	i.DB , _ = gorm.Open("postgres" , psqlInfo)
	i.DB.LogMode(true)
  	i.DB.AutoMigrate(&Update_Request{})
}

func (i *Impl) NewRequest(r Update_Request){
	i.DB.Create(&r)
}

func (i *Impl) AllRequests(r []Update_Request , ch string) ([]Update_Request){
	i.DB.Where("channel = ?" , ch).Find(&r)
	return r
}

func (i *Impl) CloseDB(){
	i.DB.Close()
}