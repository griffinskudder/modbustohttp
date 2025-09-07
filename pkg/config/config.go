package config

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"time"

	"github.com/caarlos0/env/v11"
)

// ModbusFunction is a particular Modbus function which can be configured to be supported by the application.
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

// Modbus contains Modbus protocol specific config
type Modbus struct {
	// Host is the hostname of the modbus server to connect to
	Host string `json:"host" env:"HOST" envDefault:"localhost"`
	// Port is the port of the modbus server to connect to
	Port int `json:"port" env:"PORT" envDefault:"502"`
	// SlaveID is the slave ID of the modbus server to connect to
	SlaveID byte `json:"slaveID" env:"SLAVE_ID" envDefault:"1"`
	// ConnectionTimeout is the amount of time the server will keep the connection to the modbus server active if no
	// requests are made
	ConnectionTimeout time.Duration `json:"connectionTimeout" env:"CONNECTION_TIMEOUT" envDefault:"10s"`
	// FunctionsSupported is the list of available ModbusFunction supported by the modbus server
	FunctionsSupported []ModbusFunction `json:"functionsSupported" env:"FUNCTIONS_SUPPORTED" envDefault:"ReadCoils,ReadDiscreteInputs,ReadHoldingRegisters,ReadInputRegisters,WriteSingleCoil,WriteMultipleCoils,WriteMultipleRegisters,WriteSingleRegister,MaskWriteSingleRegister"`
}

// HTTP contains the HTTP specific config of the application
type HTTP struct {
	Host string `json:"host" env:"HOST" envDefault:""`
	Port int    `json:"port" env:"PORT" envDefault:"8080"`
}

// App is the modbustohttp application config
type App struct {
	// Modbus contains modbus specific config
	Modbus Modbus `json:"modbus" envPrefix:"MODBUS_"`
	// HTTP contains HTTP specific config
	HTTP HTTP `json:"http" envPrefix:"HTTP_"`
}

// LoadAppConfig loads the application config from the given path. If path is nil then config will be loaded from
// environment variables, falling back to default values if a given environment variable is not present.
func LoadAppConfig(path *string) (*App, error) {
	var err error
	var file *os.File
	if path != nil {
		file, err = os.Open(*path)
		defer func(file *os.File) {
			_ = file.Close()
		}(file)
	}
	// Error opening the file or no file path provided
	if err != nil || path == nil {
		if errors.Is(err, fs.ErrNotExist) {
			app := App{}
			err = env.Parse(&app)
			if err != nil {
				// If there is an error when parsing the environment variables, return the error.
				return nil, err
			}
			return &app, nil
		} else {
			// If there is an error when opening the file, but the file does exist, return the error.
			return nil, err
		}
	}
	var app App
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&app)
	if err != nil {
		// If there is an error when decoding the file, return the error.
		return nil, err
	}
	return &app, nil
}
