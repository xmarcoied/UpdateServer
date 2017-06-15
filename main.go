package main

import (
    "os"
    "bufio"
    "bytes"
    "io"
    "fmt"
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
    //data := "Marco"
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

func readLines(path string) (lines []string, err error) {
    var (
        file *os.File
        part []byte
        prefix bool
    )
    if file, err = os.Open(path); err != nil {
        return
    }
    defer file.Close()

    reader := bufio.NewReader(file)
    buffer := bytes.NewBuffer(make([]byte, 0))
    for {
        if part, prefix, err = reader.ReadLine(); err != nil {
            break
        }
        buffer.Write(part)
        if !prefix {
            lines = append(lines, buffer.String())
            buffer.Reset()
        }
    }
    if err == io.EOF {
        err = nil
    }
    return
}

func writeLines(lines []string, path string) (err error) {
    var (
        file *os.File
    )

    if file, err = os.Create(path); err != nil {
        return
    }
    defer file.Close()

    //writer := bufio.NewWriter(file)
    for _,item := range lines {
        //fmt.Println(item)
        _, err := file.WriteString(strings.TrimSpace(item) + "\n"); 
        //file.Write([]byte(item)); 
        if err != nil {
            //fmt.Println("debug")
            fmt.Println(err)
            break
        }
    }
    return
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
