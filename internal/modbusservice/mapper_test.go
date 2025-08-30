package modbusservice

import (
	modbusv1alpha1 "modbustohttp/gen/modbustohttp/v1alpha1"
	"testing"
)

func TestMapByteArrayToBooleanAddress(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		startAddr   uint32
		maxQuantity uint32
		want        []*modbusv1alpha1.BooleanAddress
	}{
		{
			name:        "Single byte with max quantity 8",
			data:        []byte{0b10101010},
			startAddr:   0,
			maxQuantity: 8,
			want: []*modbusv1alpha1.BooleanAddress{
				{Address: 0, Value: false},
				{Address: 1, Value: true},
				{Address: 2, Value: false},
				{Address: 3, Value: true},
				{Address: 4, Value: false},
				{Address: 5, Value: true},
				{Address: 6, Value: false},
				{Address: 7, Value: true},
			},
		},
		{
			name:        "Two bytes with max quantity 10",
			data:        []byte{0xFF, 0x00},
			startAddr:   100,
			maxQuantity: 10,
			want: []*modbusv1alpha1.BooleanAddress{
				{Address: 100, Value: true},
				{Address: 101, Value: true},
				{Address: 102, Value: true},
				{Address: 103, Value: true},
				{Address: 104, Value: true},
				{Address: 105, Value: true},
				{Address: 106, Value: true},
				{Address: 107, Value: true},
				{Address: 108, Value: false},
				{Address: 109, Value: false},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MapByteArrayToBooleanAddress(tt.data, tt.startAddr, tt.maxQuantity)
			if len(got) != len(tt.want) {
				t.Errorf("MapByteArrayToBooleanAddress() length = %v, want %v", len(got), len(tt.want))
			}
			for i := range got {
				if got[i].Address != tt.want[i].Address || got[i].Value != tt.want[i].Value {
					t.Errorf("MapByteArrayToBooleanAddress()[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestMapByteArrayToRegisters(t *testing.T) {
	tests := []struct {
		name      string
		data      []byte
		startAddr uint32
		want      []*modbusv1alpha1.Register
	}{
		{
			name:      "Single register",
			data:      []byte{0x12, 0x34},
			startAddr: 0,
			want: []*modbusv1alpha1.Register{
				{Address: 0, Value: 0x1234},
			},
		},
		{
			name:      "Multiple registers with offset",
			data:      []byte{0x00, 0xFF, 0xFF, 0x00},
			startAddr: 100,
			want: []*modbusv1alpha1.Register{
				{Address: 100, Value: 0x00FF},
				{Address: 101, Value: 0xFF00},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MapByteArrayToRegisters(tt.data, tt.startAddr)
			if len(got) != len(tt.want) {
				t.Errorf("MapByteArrayToRegisters() length = %v, want %v", len(got), len(tt.want))
			}
			for i := range got {
				if got[i].Address != tt.want[i].Address || got[i].Value != tt.want[i].Value {
					t.Errorf("MapByteArrayToRegisters()[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestMapByteArrayToBooleanAddress_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		startAddr   uint32
		maxQuantity uint32
		want        []*modbusv1alpha1.BooleanAddress
	}{
		{
			name:        "Empty data",
			data:        []byte{},
			startAddr:   0,
			maxQuantity: 0,
			want:        []*modbusv1alpha1.BooleanAddress{},
		},
		{
			name:        "MaxQuantity less than available bits",
			data:        []byte{0xFF, 0xFF},
			startAddr:   0,
			maxQuantity: 4,
			want: []*modbusv1alpha1.BooleanAddress{
				{Address: 0, Value: true},
				{Address: 1, Value: true},
				{Address: 2, Value: true},
				{Address: 3, Value: true},
			},
		},
		{
			name:        "Large start address",
			data:        []byte{0x01},
			startAddr:   65535,
			maxQuantity: 1,
			want: []*modbusv1alpha1.BooleanAddress{
				{Address: 65535, Value: true},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MapByteArrayToBooleanAddress(tt.data, tt.startAddr, tt.maxQuantity)
			if len(got) != len(tt.want) {
				t.Errorf("length = %v, want %v", len(got), len(tt.want))
			}
			for i := range got {
				if got[i].Address != tt.want[i].Address || got[i].Value != tt.want[i].Value {
					t.Errorf("index %d = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestMapByteArrayToRegisters_EdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		data      []byte
		startAddr uint32
		want      []*modbusv1alpha1.Register
	}{
		{
			name:      "Empty data",
			data:      []byte{},
			startAddr: 0,
			want:      []*modbusv1alpha1.Register{},
		},
		{
			name:      "Max uint16 value",
			data:      []byte{0xFF, 0xFF},
			startAddr: 0,
			want: []*modbusv1alpha1.Register{
				{Address: 0, Value: 65535},
			},
		},
		{
			name:      "Multiple zero registers",
			data:      []byte{0x00, 0x00, 0x00, 0x00},
			startAddr: 1000,
			want: []*modbusv1alpha1.Register{
				{Address: 1000, Value: 0},
				{Address: 1001, Value: 0},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MapByteArrayToRegisters(tt.data, tt.startAddr)
			if len(got) != len(tt.want) {
				t.Errorf("length = %v, want %v", len(got), len(tt.want))
			}
			for i := range got {
				if got[i].Address != tt.want[i].Address || got[i].Value != tt.want[i].Value {
					t.Errorf("index %d = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}
