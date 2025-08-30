package main

import (
	"net/http"

	"github.com/goburrow/modbus"
)

type ModbusHandlerFunc func(handler *modbus.TCPClientHandler) http.HandlerFunc

var Routes = map[string]ModbusHandlerFunc{
	"/read_holding_registers/": ReadHoldingRegisters,
	"/write_single_register/":  WriteSingleRegister,
}
