package main 
import (

     	"github.com/gin-gonic/gin"
    	"updater/model"
)
// In Process : Use Gin package 
// Global variables are bad 
var db model.Impl

func main(){
		db.ConnectDB()
		defer db.CloseDB()
        router := gin.Default();

        // TODO : status generation
        vlc_router := router.Group("/vlc/:channel")
        {
        	vlc_router.StaticFile("status" , "./client/static/status")
        	vlc_router.StaticFile("status.asc" , "./client/static/status.asc")
	    	vlc_router.GET("/showoff" , showoff)
        	vlc_router.GET("/update", update)
    	}
        router.Run(":80")
}
