package main 
import (
    	"log"
    	"net/http"
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

        // TODO : Create a controller for the status managment
        router.StaticFile("status" , "./client/static/status")
        router.StaticFile("status.asc" , "./client/static/status.asc")
	    router.GET("/showoff" , showoff)
        router.POST("/update", update)
        router.Run(":80")
}

func showoff(c *gin.Context) {
    	var requests []model.Update_Request
        requests = db.AllRequests(requests)
        c.String(http.StatusOK, "Hello")
        log.Println(requests)
}

func update(c *gin.Context) {
        // Assuming VLC is sending JSON
        var request model.Update_Request
        c.BindJSON(&request)
        log.Println(request)
        // TODO : DB Model API
        // FIXME : initiate the DB once and pass it everywhere
        db.NewRequest(request)
}
