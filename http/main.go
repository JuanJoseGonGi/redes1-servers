package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := "8080"

	if len(os.Args) > 1 && os.Args[1] != "" {
		port = os.Args[1]
	}

	log.Printf("Starting at %s\n", port)

	err := http.ListenAndServe(":"+port, http.FileServer(http.Dir(".")))
	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
	}
}
