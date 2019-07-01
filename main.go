package main

import (
	"log"
	"net/http"
)


func main() {

	log.Println("Starting server...")
	services := new(Services)

	http.Handle("/", services)

	log.Println("Listening on port 80...")
	log.Fatal(http.ListenAndServe(":80", nil))
}
