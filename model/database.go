package model
import(
//	"log"
	"github.com/jinzhu/gorm"
       _"github.com/lib/pq"
)

type Update_Request struct{
	//gorm.Model
	//ID          uint   `gorm:"primary_key"`
	OS 	 		string	`json:os` 	   
	OS_VER 		string	`json:os_ver`
	OS_ARCH 	string 	`js:os_arch`
	VLC_VER 	string	`json:vlc_ver`
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
	i.DB.Create(r)
}

func (i *Impl) AllRequests(r []Update_Request) ([]Update_Request){
	i.DB.Find(&r)
	return r
}

func (i *Impl) CloseDB(){
	i.DB.Close()
}