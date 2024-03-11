package main

import (
	"demotivator-generator/api"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", api.GenerateDemotivator)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
