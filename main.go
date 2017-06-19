package main 
import (
    	"log"
    	"encoding/json"
    	"net/http"
     	//"github.com/gin-gonic/gin"
    	"updater/model"
)
// TODO : Use Gin package 
// Global variables are bad 
var db model.Impl

func main(){
		db.ConnectDB()
		defer db.CloseDB()

	    http.HandleFunc("/showoff" , showoff)
        http.HandleFunc("/update", update) 
        http.HandleFunc("/vlc/status", statusHandler)
        http.HandleFunc("/vlc/status.asc", status_ascHandler)
        log.Fatal(http.ListenAndServe(":80", nil))
}

// TODO : Create a controller for the status managment
func statusHandler(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "client/static/status")
        log.Println("status")
}
func status_ascHandler(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "client/static/status.asc")
        log.Println("status.asc")
}
func showoff(w http.ResponseWriter, r *http.Request) {
    	var requests []model.Update_Request
        requests = db.AllRequests(requests)
        log.Println(requests)
}

func update(w http.ResponseWriter, r *http.Request) {
        // Assuming VLC is sending JSON
        var request model.Update_Request
        err := json.NewDecoder(r.Body).Decode(&request)
        if err != nil {
            http.Error(w, err.Error(), 400)
            return
        }
        // TODO : DB Model API
        // FIXME : initiate the DB once and pass it everywhere
        db.NewRequest(request)
}
