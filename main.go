package main

import (

    //"fmt"
    "net/http"
    "strings"
    "log"
    "io/ioutil"
    "net/url"
    "time"
    "html/template"

)
type Page struct{
    Body []string
}
func homepage(w http.ResponseWriter, r *http.Request) {
    data, _ := readLines("in.txt") 
    t, _ := template.ParseFiles("templates/home.html")
    p := Page{Body : data}
    t.Execute(w , p)
    
}
func statusHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "static/status")
    log.Println("status")
}
func status_ascHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "static/status.asc")
    log.Println("status.asc")
}

func home(w http.ResponseWriter, r *http.Request) {
    var data []string
    tm := time.Now().Format(time.RFC3339)
    line := "Time :" + tm
    data = append(data , line)  
    if r.Method == "POST"{
        body, _ := ioutil.ReadAll(r.Body)
        values, _ := url.ParseQuery(string(body))
        for v , k := range values{
            line = v + " : " + strings.Join(k, "")
            data = append(data , line)
        }
    }
    if r.Method == "GET"{
        r.ParseForm() 
        for v , k := range r.Form {
            line = v + " : " + strings.Join(k, "")
            data = append(data , line)
        }

    }
    sdata, _ := readLines("in.txt")
    for _ , line := range data{
        sdata = append(sdata , line)
    }
    sdata = append(sdata , "--")
    _ = writeLines(sdata, "in.txt")
 
}

func main() {
    // You shouldn't use the root URL this way
    http.HandleFunc("/", home) 
    http.HandleFunc("/home", homepage) 
    http.HandleFunc("/vlc/status", statusHandler)
    http.HandleFunc("/vlc/status.asc", status_ascHandler)
    log.Fatal(http.ListenAndServe(":80", nil))
}
