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
- Write Bit In Register

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

The server can be configured using a json file. An example config file can be found [here](./config.json).

```json
{
    "modbus": {
        "protocol": "tcp",
        "host": "localhost",
        "port": 502,
        "slave_id": 1
    },
    "http": {
        "host": "localhost",
        "port": 8080
    }
}
```