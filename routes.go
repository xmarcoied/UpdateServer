package main 
import (
        "log"
        "net/http"
        "github.com/gin-gonic/gin"
        "updater/model"
        "fmt"
)
func showoff(c *gin.Context) {
    	var requests []model.Update_Request
        requests = db.AllRequests(requests)
        //fmt.Fprintln(c.Writer , requests)
        for _ , request := range requests{
        	fmt.Fprint(c.Writer , "OS :" , request.OS)
        	fmt.Fprint(c.Writer , " , OS_VER :" , request.OS_VER)
        	fmt.Fprint(c.Writer , " , OS_ARCH :" , request.OS_ARCH)
        	fmt.Fprintln(c.Writer ," , VLC_VER :" , request.VLC_VER)
        }
}

func update(c *gin.Context) {
        // Request params are now getting in GET params
        var request model.Update_Request
        c.Bind(&request)
        log.Println(request)
        // TODO : DB Model API
        // FIXME : initiate the DB once and pass it everywhere
        db.NewRequest(request)
        c.Redirect(http.StatusMovedPermanently, "http://update.videolan.org/vlc/status")

}
