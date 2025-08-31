package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type ModbusFunction string

const (
	ReadCoils               ModbusFunction = "ReadCoils"
	ReadDiscreteInputs      ModbusFunction = "ReadDiscreteInputs"
	ReadHoldingRegisters    ModbusFunction = "ReadHoldingRegisters"
	ReadInputRegisters      ModbusFunction = "ReadInputRegisters"
	WriteSingleCoil         ModbusFunction = "WriteSingleCoil"
	WriteMultipleCoils      ModbusFunction = "WriteMultipleCoils"
	WriteMultipleRegisters  ModbusFunction = "WriteMultipleRegisters"
	WriteSingleRegister     ModbusFunction = "WriteSingleRegister"
	MaskWriteSingleRegister ModbusFunction = "MaskWriteSingleRegister"
)

type Modbus struct {
	Host               string           `json:"host"`
	Port               int              `json:"port"`
	SlaveID            byte             `json:"slaveID"`
	FunctionsSupported []ModbusFunction `json:"functionsSupported"`
}

type HTTP struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type AppConfig struct {
	Modbus Modbus `json:"modbus"`
	HTTP   HTTP   `json:"http"`
}

func LoadAppConfig(path string) (*AppConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Return default configuration if file does not exist
			return &AppConfig{
				Modbus: Modbus{
					Host:    "localhost",
					Port:    502,
					SlaveID: 1,
					FunctionsSupported: []ModbusFunction{
						ReadCoils,
						ReadDiscreteInputs,
						ReadHoldingRegisters,
						ReadInputRegisters,
						WriteSingleCoil,
						WriteMultipleCoils,
						WriteMultipleRegisters,
						WriteSingleRegister,
						MaskWriteSingleRegister,
					},
				},
				HTTP: HTTP{
					Host: "localhost",
					Port: 8080,
				},
			}, nil
		} else {
			return nil, err
		}
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Error closing file:", err)
		}
	}(file)
	var app AppConfig
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&app)
	if err != nil {
		return nil, err
	}
	return &app, nil
}
