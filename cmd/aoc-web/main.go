package main

import (
	"../../internals"
	"../../internals/core"
	"log"
	"net/http"
)

func main() {
	config := core.NewConfig("config/app-config.json")

	server := &http.Server{
		Addr:    ":8080",
		Handler: internals.NewRouter(config),
	}
	log.Fatal(server.ListenAndServe())
}
