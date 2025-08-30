package main

import (
	"fmt"
	"log"
	"modbustohttp/config"
	"net/http"
	"os"
	"time"

	"github.com/goburrow/modbus"
)

type Middleware func(handler http.Handler) http.Handler

type ModbusToHTTPServer struct {
	httpConfig config.HTTP
	middleware []Middleware
	handler    *modbus.TCPClientHandler
	mux        *http.ServeMux
}

func (mts *ModbusToHTTPServer) ListenAndServe() error {
	for route, routeHandler := range Routes {
		mts.mux.HandleFunc(route, routeHandler(mts.handler))
	}
	var mux2 http.Handler = mts.mux
	for _, wrapper := range mts.middleware {
		mux2 = wrapper(mux2)
	}
	return http.ListenAndServe(fmt.Sprintf(":%d", mts.httpConfig.Port), mux2)
}

func NewModbusToHTTPServer(httpConfig config.HTTP, mux *http.ServeMux, handler *modbus.TCPClientHandler, middleware []Middleware) *ModbusToHTTPServer {
	return &ModbusToHTTPServer{
		httpConfig: httpConfig,
		middleware: middleware,
		handler:    handler,
		mux:        mux,
	}
}

func main() {
	appConfig, err := config.LoadAppConfig("config.json")
	if err != nil {
		panic(err)
	}

	handler := modbus.NewTCPClientHandler(fmt.Sprintf("%s:%d", appConfig.Modbus.Host, appConfig.Modbus.Port))
	handler.Timeout = 10 * time.Second
	handler.SlaveId = appConfig.Modbus.SlaveID
	handler.Logger = log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)
	mux := http.NewServeMux()

	middleware := []Middleware{
		LogMiddleware,
	}
	server := NewModbusToHTTPServer(appConfig.HTTP, mux, handler, middleware)
	fmt.Printf("Server starting on port %d...\n", appConfig.HTTP.Port)
	log.Fatal(server.ListenAndServe())
}
