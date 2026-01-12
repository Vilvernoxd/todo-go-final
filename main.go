package main

import (
	"log"
	"net/http"

	"todo-go-final/pkg/api"
	"todo-go-final/pkg/db"
)

func main() {
	err := db.Init("scheduler.db")
	if err != nil {
		log.Fatal(err)
	}

	api.Init()

	webDir := "./web"
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	log.Println("Server started on :7540")
	err = http.ListenAndServe(":7540", nil)
	if err != nil {
		log.Fatal(err)
	}
}
