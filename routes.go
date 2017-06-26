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
        requests = db.AllRequests(requests , c.Param("channel"))

        for _ , request := range requests{
            fmt.Fprint(c.Writer , "ID :" , request.ID)
            fmt.Fprint(c.Writer ," , created_at :" , request.Created_At)
        	fmt.Fprint(c.Writer , " , OS :" , request.OS)
        	fmt.Fprint(c.Writer , " , OS_VER :" , request.OS_VER)
        	fmt.Fprint(c.Writer , " , OS_ARCH :" , request.OS_ARCH)
        	fmt.Fprint(c.Writer ," , VLC_VER :" , request.VLC_VER)
            fmt.Fprintln(c.Writer ," , IP :" , request.IP)
        }
        
}

func update(c *gin.Context) {
        // Request params are now getting in GET params
        var request model.Update_Request
        c.Bind(&request)
        log.Println(request)
        request.Channel = c.Param("channel")
        // TODO : DB Model API
        // FIXME : initiate the DB once and pass it everywhere
        db.NewRequest(request)
        status_route := "http://update.videolan.org/vlc/" + c.Param("channel") + "/status"
        c.Redirect(http.StatusMovedPermanently, status_route)

}
