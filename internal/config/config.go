package config

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"

	"github.com/caarlos0/env/v11"
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
	Host               string           `json:"host" env:"HOST" envDefault:"localhost"`
	Port               int              `json:"port" env:"PORT" envDefault:"502"`
	SlaveID            byte             `json:"slaveID" env:"SLAVE_ID" envDefault:"1"`
	FunctionsSupported []ModbusFunction `json:"functionsSupported" env:"FUNCTIONS_SUPPORTED" envDefault:"ReadCoils,ReadDiscreteInputs,ReadHoldingRegisters,ReadInputRegisters,WriteSingleCoil,WriteMultipleCoils,WriteMultipleRegisters,WriteSingleRegister,MaskWriteSingleRegister"`
}

type HTTP struct {
	Host string `json:"host" env:"HOST" envDefault:""`
	Port int    `json:"port" env:"PORT" envDefault:"8080"`
}

type AppConfig struct {
	Modbus Modbus `json:"modbus" envPrefix:"MODBUS_"`
	HTTP   HTTP   `json:"http" envPrefix:"HTTP_"`
}

func LoadAppConfig(path *string) (*AppConfig, error) {
	file, err := os.Open(*path)
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
		}
	}(file)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			app := AppConfig{}
			err = env.Parse(&app)
			if err != nil {
				return nil, err
			}
			return &app, nil
		}
	}
	var app AppConfig
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&app)
	if err != nil {
		return nil, err
	}
	return &app, nil
}
