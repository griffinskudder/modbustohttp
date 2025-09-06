# Modbus To HTTP

A lightweight, extendable modbus to http server using [go modbus](https://pkg.go.dev/github.com/goburrow/modbus) and 
[connect-go](https://pkg.go.dev/github.com/bufbuild/connect-go).

## Supported Functions

- Read Coils
- Read Discrete Inputs
- Read Holding Registers
- Read Input Registers
- Write Single Coil
- Write Single Register
- Write Multiple Coils
- Write Multiple Registers
- Write Bit In Register (Custom Function)

### Write Bit In Register
This custom function allows you to write a single bit in a holding register without affecting the other bits. 
It is a wrapper around the MarkWriteSingleRegister function. Providing a simpler syntax which covers the use case 
of writing a single bit.

The request requires the following parameters:
- `address`: The address of the holding register (0-based).
- `bit_position`: The position of the bit to write (0-15).
- `value`: The value to write (true or false).


## Supported Modbus Protocols

- Modbus TCP

## Supported HTTP Content Types

- `application/json`
- `application/proto`
- `application/connect`

## Specs

- [OpenAPI Spec](./specs/openapi)
- [Protocol Buffers Spec](./proto/modbustohttp/v1alpha1)

## Config

The server can be configured using environment variables or a json config file. If both are provided, the json file will
take precedence. To specify a config file, use the `-config` flag when starting the server or set the `CONFIG_FILE` 
environment variable. The default config file is `config.json` in the application working directory.

### Environment Variables

The following environment variables can be used to configure the server:
- `MODBUS_HOST`: The modbus server host (default: localhost)
- `MODBUS_PORT`: The modbus server port (default: 502)
- `MODBUS_SLAVE_ID`: The modbus slave id (default: 1)
- `MODBUS_FUNCTIONS_SUPPORTED`: A comma separated list of supported modbus functions (default: all functions supported)
- `HTTP_HOST`: The http server host (default: blank, all interfaces)
- `HTTP_PORT`: The http server port (default: 8080)

### File
The server can be configured using a json file. An example config file can be found [here](config.example.json).

## Docker

A Dockerfile is provided to build a docker image of the server. To build the image, run the following command:

```bash
docker build -t modbustohttp:latest .
```

To run the server in a docker container, use the following command:

```bash
docker run -d -p 8080:8080 -e HTTP_HOST --name modbustohttp modbustohttp:latest
```

### Docker Compose
A docker compose file is provided to run the server in a docker container which can be found [here](docker-compose.yaml).

## Development
To run the server locally, you will need to have Go installed. You can then run the following command:

```bash
make run
```

This will start the server with the default configuration. You can then access the server at `http://localhost:8080`.

To run the tests, use the following command:

```bash
make test
```

To generate the protobuf and openapi files, use the following command:

```bash
make generate
```
This will generate the protobuf files in the `proto` directory and the openapi files in the `specs/openapi` directory.
To format the code, use the following command:

```bash
make fmt
``` 
This will format the code using `go fmt`.
To lint the code, use the following command:

```bash
make lint
```
This will lint the code using `golangci-lint` and `buf lint`.
To build the server, use the following command:

```bash
make build
```
This will build the server and create an executable file in the `dist` directory.
To clean the build files, use the following command:

```bash
make clean
```
This will remove the `dist` directory.