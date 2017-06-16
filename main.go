package main

import (
  	"database/sql"
  	_ "github.com/lib/pq"
    
    "updater/router"
)

func main() {
	// TODO : Move this to Model/database package
  	psqlInfo := "host=localhost dbname=updater user=postgres password=postgres sslmode=disable"
    db, _ := sql.Open("postgres", psqlInfo)
    db.Close()
    server.StartServer()
}
