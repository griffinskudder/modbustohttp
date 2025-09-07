package modbusservice

import (
	"encoding/binary"
	"modbustohttp/internal/utils"
	modbusv1alpha1 "modbustohttp/service/modbustohttp/v1alpha1"
)

func MapByteArrayToBooleanAddress(data []byte, startAddress uint32, maxQuantity uint32) []*modbusv1alpha1.BooleanAddress {
	booleanAddresses := make([]*modbusv1alpha1.BooleanAddress, maxQuantity)
	for i, dataByte := range data {
		bits := utils.ByteToBoolSlice(dataByte)
		for j, bit := range bits {
			if uint32(i*8+j) >= maxQuantity {
				break
			}
			booleanAddresses[i*8+j] = &modbusv1alpha1.BooleanAddress{
				Address: startAddress + uint32(i*8+j),
				Value:   bit,
			}
		}
	}
	return booleanAddresses
}

func MapByteArrayToRegisters(data []byte, startAddress uint32) []*modbusv1alpha1.Register {
	registers := make([]*modbusv1alpha1.Register, len(data)/2)
	for i := 0; i < len(data); i = i + 2 {
		bytes := data[i : i+2]
		registerValue := binary.BigEndian.Uint16(bytes[:2])
		registers[i/2] = &modbusv1alpha1.Register{
			Address: uint32(i/2) + startAddress,
			Value:   uint32(registerValue),
		}
	}
	return registers
}
