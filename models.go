package main

type ModbusFunction string

type ModbusRequest struct {
	function ModbusFunction
}

type ReadHoldingRegistersRequest struct {
	Address uint16 `json:"address"`
	// Default to 1 when not present or empty
	Quantity uint16 `json:"quantity,omitempty"`
}

type ReadHoldingRegistersResponse struct {
	Registers []uint16 `json:"registers"`
}

type WriteSingleRegisterRequest struct {
	Address uint16 `json:"address"`
	Value   uint16 `json:"value"`
}
