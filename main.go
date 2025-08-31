package main

import (
	"fmt"
	"log/slog"
	"modbustohttp/internal/interceptors"
	"modbustohttp/internal/modbusservice"
	"net/http"
	"os"
	"strings"
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

func setupModbusHandler(appConfig *config.AppConfig, logger *slog.Logger) *modbus.TCPClientHandler {
	logger.Info("setting up modbus handler",
		slog.String("host", appConfig.Modbus.Host),
		slog.Int("port", appConfig.Modbus.Port),
		slog.Int("slave_id", int(appConfig.Modbus.SlaveID)),
	)
	handler := modbus.NewTCPClientHandler(fmt.Sprintf("%s:%d", appConfig.Modbus.Host, appConfig.Modbus.Port))
	handler.Timeout = 10 * time.Second
	handler.SlaveId = appConfig.Modbus.SlaveID
	return handler
}

func setupReflector(mux *http.ServeMux, logger *slog.Logger) {
	names := []string{v1alpha1connect.ModbusServiceName}
	logger.Info("setting up reflector",
		slog.String("services", strings.Join(names, ",")),
	)
	reflector := grpcreflect.NewReflector(
		grpcreflect.NamerFunc(
			func() []string { return names },
		),
	)

	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))
}

func setupInterceptors(logger *slog.Logger) ([]connect.Interceptor, error) {
	logger.Info("setting up interceptors",
		slog.String("interceptors", "validation, logging"),
	)
	// Create the validation interceptor provided by connectrpc.com/validate.
	validateInterceptor, err := validate.NewInterceptor()
	if err != nil {
		logger.Error("error creating interceptor",
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	loggingInterceptor := interceptors.NewLoggingInterceptor(logger)
	return []connect.Interceptor{validateInterceptor, loggingInterceptor}, nil

}

func setupServiceHandler(
	modbusServer *modbusservice.Service,
	mux *http.ServeMux,
	logger *slog.Logger,
	serviceInterceptors ...connect.Interceptor,
) {
	logger.Info("setting up service handler",
		slog.Int("num_interceptors", len(serviceInterceptors)),
	)
	mux.Handle(v1alpha1connect.NewModbusServiceHandler(
		modbusServer,
		connect.WithInterceptors(serviceInterceptors...),
	))
}

func setupServer(addr string, mux *http.ServeMux, logger *slog.Logger) *http.Server {
	logger.Info("setting up http server",
		slog.String("addr", addr),
	)
	server := &http.Server{
		Addr:              addr,
		ReadHeaderTimeout: 3 * time.Second,
		// Use h2c so we can serve HTTP/2 without TLS.
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}
	return server
}

func main() {
	appConfig, err := config.LoadAppConfig("config.json")
	if err != nil {
		panic(err)
	}
	addr := fmt.Sprintf("%s:%d", appConfig.HTTP.Host, appConfig.HTTP.Port)

	structuredLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	handler := setupModbusHandler(appConfig, structuredLogger)
	modbusServer := modbusservice.NewService(handler)
	mux := http.NewServeMux()

	serviceInterceptors, err := setupInterceptors(structuredLogger)

	if err != nil {
		slog.Error("error setting up interceptors",
			slog.String("error", err.Error()),
		)
		return
	}

	setupServiceHandler(modbusServer, mux, structuredLogger, serviceInterceptors...)

	setupReflector(mux, structuredLogger)

	server := setupServer(addr, mux, structuredLogger)

	structuredLogger.Info("starting modbus server", slog.String("addr", addr))

	if err := server.ListenAndServe(); err != nil {
		slog.Error("error running application",
			slog.String("error", err.Error()),
		)
	}
}
