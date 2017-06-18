// TODO : Use Gin package 
package server
import (
    
   	"net/http"
    	"log"
    	"encoding/json"
    	
    	"updater/model"
)

// TODO : Create a controller for the status managment
func statusHandler(w http.ResponseWriter, r *http.Request) {
    	http.ServeFile(w, r, "client/static/status")
    	log.Println("status")
}
func status_ascHandler(w http.ResponseWriter, r *http.Request) {
    	http.ServeFile(w, r, "client/static/status.asc")
    	log.Println("status.asc")
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
        db := model.Impl{}
        db.ConnectDB()
        db.NewRequest(request)        
}

func StartServer() {

    	http.HandleFunc("/update", update) 
    	http.HandleFunc("/vlc/status", statusHandler)
    	http.HandleFunc("/vlc/status.asc", status_ascHandler)
    	log.Fatal(http.ListenAndServe(":80", nil))
}