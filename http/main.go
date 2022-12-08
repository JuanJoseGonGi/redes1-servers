package main

import (
	"log"
	"net/http"
	"strconv"
)

func main() {
	port := 8080
	log.Printf("Starting at %d\n", port)

	err := http.ListenAndServe(":"+strconv.Itoa(port), http.FileServer(http.Dir(".")))
	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
	}
}
