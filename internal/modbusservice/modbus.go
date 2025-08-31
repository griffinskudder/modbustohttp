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

	registers := MapByteArrayToRegisters(modbusData, req.Msg.GetAddress())
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
	coils := MapByteArrayToBooleanAddress(data, req.Msg.GetAddress(), req.Msg.GetQuantity())
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

func (s Service) MaskWriteRegister(
	_ context.Context,
	req *connect.Request[modbusv1alpha1.MaskWriteRegisterRequest],
) (*connect.Response[modbusv1alpha1.MaskWriteRegisterResponse], error) {
	err := s.modbusHandler.Connect()
	if err != nil {
		return nil, err
	}
	client := modbus.NewClient(s.modbusHandler)
	var andMask uint16
	if req.Msg.GetBitAndMask() != nil {
		andMaskBytes := utils.BoolArrayToByteArray(req.Msg.GetBitAndMask().GetBits())
		binary.BigEndian.Uint16(andMaskBytes[:2])
	} else {
		andMask = uint16(req.Msg.GetIntAndMask())
	}
	var orMask uint16
	if req.Msg.GetBitOrMask() != nil {
		orMaskBytes := utils.BoolArrayToByteArray(req.Msg.GetBitOrMask().GetBits())
		binary.BigEndian.Uint16(orMaskBytes[:2])
	} else {
		orMask = uint16(req.Msg.GetIntOrMask())
	}

	_, err = client.MaskWriteRegister(
		uint16(req.Msg.GetAddress()),
		andMask,
		orMask,
	)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&modbusv1alpha1.MaskWriteRegisterResponse{}), nil
}

func NewService(modbusHandler *modbus.TCPClientHandler) *Service {
	return &Service{
		modbusHandler,
	}
}
