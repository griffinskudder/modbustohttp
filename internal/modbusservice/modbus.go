package modbusservice

import (
	"context"
	"encoding/binary"
	"modbustohttp/internal/config"
	"modbustohttp/internal/utils"
	"slices"

	"connectrpc.com/connect"
	"github.com/goburrow/modbus"

	modbusv1alpha1 "modbustohttp/service/modbustohttp/v1alpha1"
)

type Service struct {
	modbusHandler *modbus.TCPClientHandler
	modbusConfig  *config.Modbus
}

func (s Service) ReadHoldingRegisters(
	_ context.Context,
	req *connect.Request[modbusv1alpha1.ReadHoldingRegistersRequest],
) (*connect.Response[modbusv1alpha1.ReadHoldingRegistersResponse], error) {
	if slices.Index(s.modbusConfig.FunctionsSupported, config.ReadHoldingRegisters) == -1 {
		return nil, connect.NewError(connect.CodeUnimplemented, nil)
	}
	err := s.modbusHandler.Connect()
	if err != nil {
		return nil, err
	}
	client := modbus.NewClient(s.modbusHandler)
	modbusData, err := client.ReadHoldingRegisters(
		uint16(req.Msg.GetAddress()),
		uint16(req.Msg.GetQuantity()),
	)
	if err != nil {
		return nil, err
	}

	registers := MapByteArrayToRegisters(modbusData, req.Msg.GetAddress())
	return connect.NewResponse(&modbusv1alpha1.ReadHoldingRegistersResponse{Registers: registers}), nil
}

func (s Service) WriteSingleRegister(
	_ context.Context,
	req *connect.Request[modbusv1alpha1.WriteSingleRegisterRequest],
) (*connect.Response[modbusv1alpha1.WriteSingleRegisterResponse], error) {
	if slices.Index(s.modbusConfig.FunctionsSupported, config.WriteSingleRegister) == -1 {
		return nil, connect.NewError(connect.CodeUnimplemented, nil)
	}
	err := s.modbusHandler.Connect()
	if err != nil {
		return nil, err
	}
	client := modbus.NewClient(s.modbusHandler)

	_, err = client.WriteSingleRegister(uint16(req.Msg.GetRegister().Address), uint16(req.Msg.GetRegister().Value))
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&modbusv1alpha1.WriteSingleRegisterResponse{}), nil
}

func (s Service) ReadCoils(
	_ context.Context,
	req *connect.Request[modbusv1alpha1.ReadCoilsRequest],
) (*connect.Response[modbusv1alpha1.ReadCoilsResponse], error) {
	if slices.Index(s.modbusConfig.FunctionsSupported, config.ReadCoils) == -1 {
		return nil, connect.NewError(connect.CodeUnimplemented, nil)
	}
	err := s.modbusHandler.Connect()
	if err != nil {
		return nil, err
	}
	client := modbus.NewClient(s.modbusHandler)

	data, err := client.ReadCoils(uint16(req.Msg.GetAddress()), uint16(req.Msg.GetQuantity()))
	if err != nil {
		return nil, err
	}
	coils := MapByteArrayToBooleanAddress(data, req.Msg.GetAddress(), req.Msg.GetQuantity())
	response := modbusv1alpha1.ReadCoilsResponse{Coils: coils}
	return connect.NewResponse(&response), nil
}

func (s Service) ReadDiscreteInputs(
	_ context.Context,
	req *connect.Request[modbusv1alpha1.ReadDiscreteInputsRequest],
) (*connect.Response[modbusv1alpha1.ReadDiscreteInputsResponse], error) {
	if slices.Index(s.modbusConfig.FunctionsSupported, config.ReadDiscreteInputs) == -1 {
		return nil, connect.NewError(connect.CodeUnimplemented, nil)
	}
	err := s.modbusHandler.Connect()
	if err != nil {
		return nil, err
	}
	client := modbus.NewClient(s.modbusHandler)

	data, err := client.ReadDiscreteInputs(uint16(req.Msg.GetAddress()), uint16(req.Msg.GetQuantity()))
	if err != nil {
		return nil, err
	}

	discreteInputs := MapByteArrayToBooleanAddress(data, req.Msg.GetAddress(), req.Msg.GetQuantity())
	response := modbusv1alpha1.ReadDiscreteInputsResponse{
		Inputs: discreteInputs,
	}
	return connect.NewResponse(&response), nil
}

func (s Service) WriteSingleCoil(
	_ context.Context,
	req *connect.Request[modbusv1alpha1.WriteSingleCoilRequest],
) (*connect.Response[modbusv1alpha1.WriteSingleCoilResponse], error) {
	if slices.Index(s.modbusConfig.FunctionsSupported, config.WriteSingleCoil) == -1 {
		return nil, connect.NewError(connect.CodeUnimplemented, nil)
	}
	err := s.modbusHandler.Connect()
	if err != nil {
		return nil, err
	}
	client := modbus.NewClient(s.modbusHandler)
	var value uint16
	switch req.Msg.GetCoil().Value {
	case true:
		value = 0xFF00
	case false:
		value = 0x0000
	}

	_, err = client.WriteSingleCoil(
		uint16(req.Msg.GetCoil().Address),
		value,
	)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&modbusv1alpha1.WriteSingleCoilResponse{}), nil
}

func (s Service) WriteMultipleCoils(
	_ context.Context,
	req *connect.Request[modbusv1alpha1.WriteMultipleCoilsRequest],
) (*connect.Response[modbusv1alpha1.WriteMultipleCoilsResponse], error) {
	if slices.Index(s.modbusConfig.FunctionsSupported, config.WriteMultipleCoils) == -1 {
		return nil, connect.NewError(connect.CodeUnimplemented, nil)
	}
	err := s.modbusHandler.Connect()
	if err != nil {
		return nil, err
	}
	client := modbus.NewClient(s.modbusHandler)
	data := utils.BoolArrayToByteArray(req.Msg.GetValues())
	_, err = client.WriteMultipleCoils(uint16(req.Msg.GetAddress()), uint16(len(req.Msg.GetValues())), data)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&modbusv1alpha1.WriteMultipleCoilsResponse{}), nil
}

func (s Service) ReadInputRegisters(
	_ context.Context,
	req *connect.Request[modbusv1alpha1.ReadInputRegistersRequest],
) (*connect.Response[modbusv1alpha1.ReadInputRegistersResponse], error) {
	if slices.Index(s.modbusConfig.FunctionsSupported, config.ReadInputRegisters) == -1 {
		return nil, connect.NewError(connect.CodeUnimplemented, nil)
	}
	err := s.modbusHandler.Connect()
	if err != nil {
		return nil, err
	}
	client := modbus.NewClient(s.modbusHandler)
	modbusData, err := client.ReadInputRegisters(
		uint16(req.Msg.GetAddress()),
		uint16(req.Msg.GetQuantity()),
	)
	if err != nil {
		return nil, err
	}

	registers := MapByteArrayToRegisters(modbusData, req.Msg.GetAddress())
	return connect.NewResponse(&modbusv1alpha1.ReadInputRegistersResponse{Registers: registers}), nil
}

func (s Service) WriteMultipleRegisters(
	_ context.Context,
	req *connect.Request[modbusv1alpha1.WriteMultipleRegistersRequest],
) (*connect.Response[modbusv1alpha1.WriteMultipleRegistersResponse], error) {
	if slices.Index(s.modbusConfig.FunctionsSupported, config.WriteMultipleRegisters) == -1 {
		return nil, connect.NewError(connect.CodeUnimplemented, nil)
	}
	err := s.modbusHandler.Connect()
	if err != nil {
		return nil, err
	}
	client := modbus.NewClient(s.modbusHandler)
	data := make([]byte, len(req.Msg.GetValues())*2)
	for i, value := range req.Msg.GetValues() {
		binary.BigEndian.PutUint16(data[i*2:i*2+2], uint16(value))
	}
	_, err = client.WriteMultipleRegisters(uint16(req.Msg.GetAddress()), uint16(len(req.Msg.GetValues())), data)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&modbusv1alpha1.WriteMultipleRegistersResponse{}), nil
}

func (s Service) WriteBitInRegister(
	_ context.Context,
	req *connect.Request[modbusv1alpha1.WriteBitInRegisterRequest],
) (*connect.Response[modbusv1alpha1.WriteBitInRegisterResponse], error) {
	primaryEnabled := slices.Index(s.modbusConfig.FunctionsSupported, config.MaskWriteSingleRegister) >= 0
	// Fallback is only possible if both WriteSingleRegister and ReadHoldingRegisters are supported
	// as we need to read the current value of the register, modify the specific bit and write it back.
	// If either of these functions is not supported, we cannot use the fallback method.
	fallbackEnabled := slices.Index(s.modbusConfig.FunctionsSupported, config.WriteSingleRegister) >= 0 && slices.Index(s.modbusConfig.FunctionsSupported, config.ReadHoldingRegisters) >= 0
	// If neither primary nor fallback method is possible, return unimplemented error.
	if !primaryEnabled && !fallbackEnabled {
		return nil, connect.NewError(connect.CodeUnimplemented, nil)
	}
	err := s.modbusHandler.Connect()
	if err != nil {
		return nil, err
	}
	client := modbus.NewClient(s.modbusHandler)

	if primaryEnabled {
		// Use MaskWriteSingleRegister if supported as it is atomic and therefore will not lead to race conditions.
		// Read the current value of the register is not necessary as MaskWriteSingleRegister allows to modify
		// specific bits directly.

		var orMask uint16
		var andMask uint16
		if req.Msg.GetValue() {
			andMask = ^uint16(0)           // All bits set to 1
			orMask = 1 << req.Msg.GetBit() // Only the target bit set to 1
		} else {
			andMask = ^(1 << req.Msg.GetBit()) // Only the target bit set to 0
			// orMask has all bits set to 0 by default
		}

		_, err = client.MaskWriteRegister(uint16(req.Msg.GetAddress()), andMask, orMask)
		if err != nil {
			return nil, err
		}
	} else {
		// Fallback if MaskWriteSingleRegister is not supported. Will read the register, modify the bit and write it
		// back. Warning: This is not atomic and can lead to race conditions if multiple clients are writing to the same
		// register.

		// Read the current value of the register
		currentData, err := client.ReadHoldingRegisters(uint16(req.Msg.GetAddress()), 1)
		if err != nil {
			return nil, err
		}
		currentValue := binary.BigEndian.Uint16(currentData)

		// Modify the specific bit
		var newValue uint16
		if req.Msg.GetValue() {
			newValue = currentValue | (1 << req.Msg.GetBit())
		} else {
			newValue = currentValue & ^(1 << req.Msg.GetBit())
		}

		// Write the new value back to the register
		_, err = client.WriteSingleRegister(uint16(req.Msg.GetAddress()), newValue)
		if err != nil {
			return nil, err
		}
	}

	return connect.NewResponse(&modbusv1alpha1.WriteBitInRegisterResponse{}), nil
}

func NewService(modbusHandler *modbus.TCPClientHandler, modbusConfig *config.Modbus) *Service {
	return &Service{
		modbusHandler,
		modbusConfig,
	}
}
