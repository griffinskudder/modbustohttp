package main

import (
	"fmt"
	"log"
	"log/slog"
	"modbustohttp/internal/modbusservice"
	"net/http"
	"os"
	"time"

	"modbustohttp/config"
	"modbustohttp/gen/modbustohttp/v1alpha1/v1alpha1connect"

	"connectrpc.com/connect"
	"connectrpc.com/grpcreflect"
	"connectrpc.com/validate"
	"github.com/goburrow/modbus"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	addr := "localhost:8080"
	appConfig, err := config.LoadAppConfig("config.json")
	if err != nil {
		panic(err)
	}

	handler := modbus.NewTCPClientHandler(fmt.Sprintf("%s:%d", appConfig.Modbus.Host, appConfig.Modbus.Port))
	handler.Timeout = 10 * time.Second
	handler.SlaveId = appConfig.Modbus.SlaveID
	handler.Logger = log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)
	modbusServer := modbusservice.NewService(handler)
	mux := http.NewServeMux()

	// Create the validation interceptor provided by connectrpc.com/validate.
	validateInterceptor, err := validate.NewInterceptor()
	if err != nil {
		slog.Error("error creating interceptor",
			slog.String("error", err.Error()),
		)
		return
	}

	mux.Handle(v1alpha1connect.NewModbusServiceHandler(
		modbusServer,
		connect.WithInterceptors(validateInterceptor),
	))

	names := []string{v1alpha1connect.ModbusServiceName}

	reflector := grpcreflect.NewReflector(
		grpcreflect.NamerFunc(
			func() []string { return names },
		),
	)

	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	slog.Info("starting modbus server", slog.String("addr", addr))
	server := &http.Server{
		Addr:              addr,
		ReadHeaderTimeout: 3 * time.Second,
		// Use h2c so we can serve HTTP/2 without TLS.
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}

	if err := server.ListenAndServe(); err != nil {
		slog.Error("error running application",
			slog.String("error", err.Error()),
		)
	}
}
