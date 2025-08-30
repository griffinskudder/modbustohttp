package modbusservice

import (
	"context"
	"encoding/binary"
	"modbustohttp/internal/utils"

	"connectrpc.com/connect"
	"github.com/goburrow/modbus"

	modbusv1alpha1 "modbustohttp/gen/modbustohttp/v1alpha1"
)

type Service struct {
	modbusHandler *modbus.TCPClientHandler
}

func (s Service) ReadHoldingRegisters(
	_ context.Context,
	req *connect.Request[modbusv1alpha1.ReadHoldingRegistersRequest],
) (*connect.Response[modbusv1alpha1.ReadHoldingRegistersResponse], error) {
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

	registers := make([]*modbusv1alpha1.Register, len(modbusData)/2)
	for i := 0; i < len(modbusData); i = i + 2 {
		bytes := modbusData[i : i+2]
		registerValue := binary.BigEndian.Uint16(bytes[:2])
		registers[i/2] = &modbusv1alpha1.Register{
			Address: uint32(i/2) + req.Msg.GetAddress(),
			Value:   uint32(registerValue),
		}
	}
	return connect.NewResponse(&modbusv1alpha1.ReadHoldingRegistersResponse{Registers: registers}), nil
}

func (s Service) WriteSingleRegister(
	_ context.Context,
	req *connect.Request[modbusv1alpha1.WriteSingleRegisterRequest],
) (*connect.Response[modbusv1alpha1.WriteSingleRegisterResponse], error) {
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
	err := s.modbusHandler.Connect()
	if err != nil {
		return nil, err
	}
	client := modbus.NewClient(s.modbusHandler)

	data, err := client.ReadCoils(uint16(req.Msg.GetAddress()), uint16(req.Msg.GetQuantity()))
	if err != nil {
		return nil, err
	}
	coils := make([]*modbusv1alpha1.BooleanAddress, req.Msg.GetQuantity())
	for i, dataByte := range data {
		bits := utils.ByteToBoolArray(dataByte)
		for j, bit := range bits {
			if req.Msg.GetAddress()+uint32(i+j) > req.Msg.GetQuantity() {
				break
			}
			coils[i+j] = &modbusv1alpha1.BooleanAddress{
				Address: req.Msg.GetAddress() + uint32(i+j),
				Value:   bit,
			}
		}
	}
	response := modbusv1alpha1.ReadCoilsResponse{Coils: coils}
	return connect.NewResponse(&response), nil
}

func (s Service) ReadDiscreteInputs(
	_ context.Context,
	req *connect.Request[modbusv1alpha1.ReadDiscreteInputsRequest],
) (*connect.Response[modbusv1alpha1.ReadDiscreteInputsResponse], error) {
	err := s.modbusHandler.Connect()
	if err != nil {
		return nil, err
	}
	client := modbus.NewClient(s.modbusHandler)

	data, err := client.ReadDiscreteInputs(uint16(req.Msg.GetAddress()), uint16(req.Msg.GetQuantity()))
	if err != nil {
		return nil, err
	}

	discreteInputs := make([]*modbusv1alpha1.BooleanAddress, req.Msg.GetQuantity())
	for i, dataByte := range data {
		bits := utils.ByteToBoolArray(dataByte)
		for j, bit := range bits {
			if req.Msg.GetAddress()+uint32(i+j) > req.Msg.GetQuantity() {
				break
			}
			discreteInputs[i+j] = &modbusv1alpha1.BooleanAddress{
				Address: req.Msg.GetAddress() + uint32(i+j),
				Value:   bit,
			}
		}
	}
	response := modbusv1alpha1.ReadDiscreteInputsResponse{
		Inputs: discreteInputs,
	}
	return connect.NewResponse(&response), nil
}

func NewService(modbusHandler *modbus.TCPClientHandler) *Service {
	return &Service{
		modbusHandler,
	}
}
