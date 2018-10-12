package main

import (
	"../../internals"
	"log"
	"net/http"
)

func main() {
	server := &http.Server{
		Addr:    ":8080",
		Handler: internals.NewRouter("config/routes"),
	}
	log.Fatal(server.ListenAndServe())
}
