package main

import (
	"github.com/gorilla/mux"
	"github.com/margostino/lagom/api"
	"github.com/margostino/lagom/configuration"
	"github.com/margostino/lagom/loader"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

func main() {
	var config = configuration.GetConfig("./config.yml")
	loadgen := loader.NewLoadGen(config)
	go loadgen.Run()
	buildRouter()
}

func buildRouter() {
	router := mux.NewRouter()
	router.Path("/monitoring").Handler(promhttp.Handler()).Methods("GET")
	router.HandleFunc("/configuration", api.UpdateConfiguration).Methods("POST")
	err := http.ListenAndServe(":9000", router)
	log.Fatal(err)
}
