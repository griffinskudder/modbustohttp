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

The request requires the following parameters:
- `address`: The address of the holding register (0-based).
- `bit_position`: The position of the bit to write (0-15).
- `value`: The value to write (true or false).

The server will read the current value of the holding register, modify the specified bit, and write the new value back 
to the register. Due to this read-modify-write operation, this function is not atomic and may lead to race conditions if
multiple clients attempt to write to the same register simultaneously.

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

The server can be configured using a json file. An example config file can be found [here](config.example.json).

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