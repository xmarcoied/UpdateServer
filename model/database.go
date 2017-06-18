package model
import(
	"log"

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
	psqlInfo := "host=localhost dbname=updater user=postgres password=postgres sslmode=disable"
	i.DB , _ = gorm.Open("postgres" , psqlInfo)
  	i.DB.AutoMigrate(&Update_Request{})
  	i.DB.SingularTable(true)
  	i.DB.LogMode(true)
  	defer i.DB.Close()
  	log.Println("Done with db")
}

func (i *Impl) NewRequest(r Update_Request){
	i.DB.Create(r)
	log.Println(r)	
}