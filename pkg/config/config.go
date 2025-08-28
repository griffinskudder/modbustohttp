package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Modbus struct {
	Host    string `json:"host"`
	Port    int    `json:"port"`
	SlaveID int    `json:"slaveID"`
}

type HTTP struct {
	Port int `json:"port"`
}

type App struct {
	Modbus Modbus `json:"modbus"`
	HTTP   HTTP   `json:"http"`
}

func LoadAppConfig(path string) (*App, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Error closing file:", err)
		}
	}(file)
	var app App
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&app)
	if err != nil {
		return nil, err
	}
	return &app, nil
}
