package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Modbus struct {
	Host    string `json:"host"`
	Port    int    `json:"port"`
	SlaveID byte   `json:"slaveID"`
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
