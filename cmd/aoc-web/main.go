package main

import (
	"../../internals"
	"../../internals/core"
	"../../internals/service"
	"log"
	"net/http"
)

func main() {
	config := core.NewConfig("config/app-config.json")
	registerServices()

	server := &http.Server{
		Addr:    ":8080",
		Handler: internals.NewRouter(config),
	}
	log.Fatal(server.ListenAndServe())
}

func registerServices() {
	core.ServicesList = map[string]interface{} {
		"Index": (*service.IndexService)(nil),
		"Solver": (*service.SolverService)(nil),
		"YearIndex": (*service.YearIndexService)(nil),
	}
}
