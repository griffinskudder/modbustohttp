package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/goburrow/modbus"
)

func main() {
	appConfig, err := LoadAppConfig("config.json")
	if err != nil {
		panic(err)
	}
	handler := modbus.NewTCPClientHandler(fmt.Sprintf("%s:%d", appConfig.Modbus.Host, appConfig.Modbus.Port))
	handler.Timeout = 10 * time.Second
	handler.SlaveId = appConfig.Modbus.SlaveID
	handler.Logger = log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)
	err = handler.Connect()
	if err != nil {
		panic(err)
	}
	defer func(handler *modbus.TCPClientHandler) {
		err := handler.Close()
		if err != nil {
			panic(err)
		}
	}(handler)
	mux := http.NewServeMux()
	for route, routeHandler := range Routes {
		mux.HandleFunc(route, routeHandler(handler))
	}
	fmt.Printf("Server starting on port %d...\n", appConfig.HTTP.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", appConfig.HTTP.Port), mux))
}
