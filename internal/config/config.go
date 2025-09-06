package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
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
	Host               string           `json:"host" env:"HOST"`
	Port               int              `json:"port" env:"PORT"`
	SlaveID            byte             `json:"slaveID" env:"SLAVE_ID"`
	FunctionsSupported []ModbusFunction `json:"functionsSupported" env:"FUNCTIONS_SUPPORTED"`
}

type HTTP struct {
	Host string `json:"host" env:"HOST"`
	Port int    `json:"port" env:"PORT"`
}

type AppConfig struct {
	Modbus Modbus `json:"modbus" env:"MODBUS_"`
	HTTP   HTTP   `json:"http" env:"HTTP_"`
}

func LoadAppConfig(path string) (*AppConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Return default configuration if file does not exist
			modbusHost, found := os.LookupEnv("MODBUS_HOST")
			if !found {
				modbusHost = "localhost"
			}
			modbusPortEnv := os.Getenv("MODBUS_PORT")
			if modbusPortEnv == "" {
				modbusPortEnv = "502"
			}
			modbusPort, err := strconv.Atoi(modbusPortEnv)
			if err != nil {
				modbusPort = 502
			}
			modbusSlaveID := os.Getenv("MODBUS_SLAVE_ID")
			if modbusSlaveID == "" {
				modbusSlaveID = "1"
			}
			slaveID, err := strconv.Atoi(modbusSlaveID)
			if err != nil || slaveID < 0 || slaveID > 255 {
				slaveID = 1
			}
			httpHost, found := os.LookupEnv("HTTP_HOST")
			if !found {
				httpHost = "localhost"
			}
			httpPortEnv := os.Getenv("HTTP_PORT")
			if httpPortEnv == "" {
				httpPortEnv = "8080"
			}
			httpPort, err := strconv.Atoi(httpPortEnv)
			if err != nil {
				httpPort = 8080
			}
			fmt.Println("Configuration file not found, using default configuration")
			return &AppConfig{
				Modbus: Modbus{
					Host:    modbusHost,
					Port:    modbusPort,
					SlaveID: byte(slaveID),
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
					Host: httpHost,
					Port: httpPort,
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
