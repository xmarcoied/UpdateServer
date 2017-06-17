package main

import (
	"github.com/jinzhu/gorm"
       _"github.com/lib/pq"
    
	"updater/router"
)
type Update_Request struct{
	gorm.Model
	ID          	uint   `gorm:"primary_key"`
	Os 	 	string 	   
	Os_ver 		string
	Os_arch 	string 
	Vlc_ver 	string
}

func main() {
	// TODO : Move this to Model/database package
	// TODO : Move the psqlinfo to config
	psqlInfo := "host=localhost dbname=updater user=postgres password=postgres sslmode=disable"
	db , _ := gorm.Open("postgres" , psqlInfo)
	defer db.Close()
  	db.AutoMigrate(&Update_Request{})
  	db.LogMode(true)

    	server.StartServer()
}
