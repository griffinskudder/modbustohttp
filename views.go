package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/goburrow/modbus"
)

func Hello(handler *modbus.TCPClientHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		client := modbus.NewClient(handler)
		err := handler.Connect()
		if err != nil {
			return
		}
		defer func(handler *modbus.TCPClientHandler) {
			err := handler.Close()
			if err != nil {
				return
			}
		}(handler)
		registers, err := client.ReadHoldingRegisters(1, 5)
		_, err = fmt.Fprintln(w, registers)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte(err.Error()))
			if err != nil {
				return
			}
		}
	}
}

// ReadHoldingRegisters reads the contents of a contiguous block of
// holding registers in the configured remove device and returns register values.
func ReadHoldingRegisters(handler *modbus.TCPClientHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		client := modbus.NewClient(handler)
		err := handler.Connect()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte(err.Error()))
			return
		}
		defer func(handler *modbus.TCPClientHandler) {
			err := handler.Close()
			if err != nil {
				return
			}
		}(handler)
		parser := json.NewDecoder(request.Body)
		requestData := ReadHoldingRegistersRequest{}
		err = parser.Decode(&requestData)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte(err.Error()))
			return
		}
		quantity := requestData.Quantity
		if quantity == 0 {
			quantity = 1
		}
		modbusData, err := client.ReadHoldingRegisters(requestData.Address, quantity)
		registers := make([]uint16, len(modbusData)/2)
		for i := 0; i < len(modbusData); i = i + 2 {
			bytes := modbusData[i : i+2]
			registerValue := binary.BigEndian.Uint16(bytes[:2])
			registers[i/2] = registerValue
		}

		response := ReadHoldingRegistersResponse{registers}

		encoder := json.NewEncoder(w)
		err = encoder.Encode(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte(err.Error()))
			if err != nil {
				return
			}
		}
	}
}

// WriteSingleRegister writes a single holding register in a remote
// device and returns register value.
func WriteSingleRegister(handler *modbus.TCPClientHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		client := modbus.NewClient(handler)
		err := handler.Connect()
		if err != nil {
			return
		}
		defer func(handler *modbus.TCPClientHandler) {
			err := handler.Close()
			if err != nil {
				return
			}
		}(handler)
		parser := json.NewDecoder(request.Body)
		requestData := WriteSingleRegisterRequest{}
		err = parser.Decode(&requestData)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte(err.Error()))
			if err != nil {
				panic(err)
			}
		}
		_, err = client.WriteSingleRegister(requestData.Address, requestData.Value)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte(err.Error()))
			if err != nil {
				panic(err)
			}
		}
		w.WriteHeader(http.StatusOK)
	}
}
