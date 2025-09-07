package main

import (
	"flag"
	"fmt"
	"log/slog"
	"modbustohttp/internal/config"
	"modbustohttp/internal/interceptors"
	"modbustohttp/internal/modbusservice"
	"net/http"
	"os"
	"strings"
	"time"

	"modbustohttp/service/modbustohttp/v1alpha1/v1alpha1connect"

	"connectrpc.com/connect"
	"connectrpc.com/grpchealth"
	"connectrpc.com/grpcreflect"
	"connectrpc.com/validate"
	"github.com/goburrow/modbus"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func setupModbusHandler(modbusConfig *config.Modbus, logger *slog.Logger) *modbus.TCPClientHandler {
	logger.Info("setting up modbus handler",
		slog.String("host", modbusConfig.Host),
		slog.Int("port", modbusConfig.Port),
		slog.Int("slave_id", int(modbusConfig.SlaveID)),
		slog.Duration("connection_timeout", modbusConfig.ConnectionTimeout),
		slog.String("functions_supported", strings.Join(func() []string {
			var funcs []string
			for _, f := range modbusConfig.FunctionsSupported {
				funcs = append(funcs, string(f))
			}
			return funcs
		}(), ", ")),
	)
	handler := modbus.NewTCPClientHandler(fmt.Sprintf("%s:%d", modbusConfig.Host, modbusConfig.Port))
	handler.Timeout = modbusConfig.ConnectionTimeout
	handler.SlaveId = modbusConfig.SlaveID
	return handler
}

func setupReflector(mux *http.ServeMux, logger *slog.Logger) {
	names := []string{v1alpha1connect.ModbusServiceName, "grpc.health.v1.Health"}
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

func setupHealthCheck(mux *http.ServeMux, logger *slog.Logger) {
	logger.Info("setting up health check")
	checker := grpchealth.NewStaticChecker(
		"modbustohttp.v1alpha1.ModbusService",
	)
	mux.Handle(grpchealth.NewHandler(checker))
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

	structuredLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	configLocation := flag.String("config", os.Getenv("CONFIG_FILE"), "location of config file")
	flag.Parse()
	if *configLocation == "" {
		*configLocation = "config.json"
	}
	appConfig, err := config.LoadAppConfig(configLocation)
	if err != nil {
		panic(err)
	}
	structuredLogger.Info(
		"loaded application configuration",
		slog.String("config_file", *configLocation),
		slog.Any("config", appConfig),
	)
	addr := fmt.Sprintf("%s:%d", appConfig.HTTP.Host, appConfig.HTTP.Port)

	handler := setupModbusHandler(&appConfig.Modbus, structuredLogger)
	modbusServer := modbusservice.NewService(handler, &appConfig.Modbus)
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

	setupHealthCheck(mux, structuredLogger)

	server := setupServer(addr, mux, structuredLogger)

	structuredLogger.Info("starting http server", slog.String("addr", addr))

	if err := server.ListenAndServe(); err != nil {
		slog.Error("error running application",
			slog.String("error", err.Error()),
		)
	}
}
