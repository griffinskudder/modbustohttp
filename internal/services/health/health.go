package health

import (
	"context"

	"connectrpc.com/grpchealth"
	"github.com/goburrow/modbus"
)

type ModbusChecker struct {
	ModbusHandler *modbus.TCPClientHandler
}

func (m ModbusChecker) Check(_ context.Context, _ *grpchealth.CheckRequest) (*grpchealth.CheckResponse, error) {
	err := m.ModbusHandler.Connect()
	if err != nil {
		return &grpchealth.CheckResponse{Status: grpchealth.StatusNotServing}, nil
	} else {
		return &grpchealth.CheckResponse{Status: grpchealth.StatusServing}, nil
	}
}

func NewModbusChecker(modbusHandler *modbus.TCPClientHandler) grpchealth.Checker {
	return ModbusChecker{ModbusHandler: modbusHandler}
}
