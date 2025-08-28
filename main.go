package main

import (
	"fmt"
	"log"
	"modbustohttp/pkg/config"
	"modbustohttp/pkg/routes"
	"net/http"
)

func main() {
	appConfig, err := config.LoadAppConfig("config.json")
	if err != nil {
		panic(err)
	}
	mux := http.NewServeMux()
	for route, handler := range routes.Routes {
		mux.HandleFunc(route, handler)
	}
	fmt.Printf("Server starting on port %d...\n", appConfig.HTTP.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", appConfig.HTTP.Port), mux))
}
