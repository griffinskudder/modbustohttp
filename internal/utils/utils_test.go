package utils

import (
	"fmt"
	"reflect"
	"testing"
)

func TestByteToBoolArray(t *testing.T) {
	tests := []struct {
		name     string
		input    byte
		expected [8]bool
	}{
		{
			name:     "Zero byte",
			input:    0x00,
			expected: [8]bool{false, false, false, false, false, false, false, false},
		},
		{
			name:     "All ones",
			input:    0xFF,
			expected: [8]bool{true, true, true, true, true, true, true, true},
		},
		{
			name:     "Alternating bits",
			input:    0x55, // 01010101
			expected: [8]bool{true, false, true, false, true, false, true, false},
		},
		{
			name:     "Single bit",
			input:    0x08, // 00001000
			expected: [8]bool{false, false, false, true, false, false, false, false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ByteToBoolSlice(tt.input)
			for i, item := range result {
				if item != tt.expected[i] {
					fmt.Printf("Index: %d, Got: %v, Expected: %v\n", i, item, tt.expected[i])
					t.Errorf("ByteToBoolSlice() = %v, want %v", result, tt.expected)
					break
				}
			}
		})
	}
}

func TestBitsToByte(t *testing.T) {
	tests := []struct {
		name     string
		input    [8]bool
		expected byte
	}{
		{
			name:     "All zeros",
			input:    [8]bool{false, false, false, false, false, false, false, false},
			expected: 0x00,
		},
		{
			name:     "All ones",
			input:    [8]bool{true, true, true, true, true, true, true, true},
			expected: 0xFF,
		},
		{
			name:     "Alternating bits",
			input:    [8]bool{true, false, true, false, true, false, true, false},
			expected: 0x55,
		},
		{
			name:     "Single bit",
			input:    [8]bool{false, false, false, true, false, false, false, false},
			expected: 0x08,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := bitsToByte(tt.input)
			if result != tt.expected {
				t.Errorf("bitsToByte() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestBoolArrayToByteArray(t *testing.T) {
	tests := []struct {
		name     string
		input    []bool
		expected []byte
	}{
		{
			name:     "Empty array",
			input:    []bool{},
			expected: []byte{},
		},
		{
			name:     "Single byte aligned",
			input:    []bool{true, false, true, false, true, false, true, false},
			expected: []byte{0x55},
		},
		{
			name:     "Multiple bytes aligned",
			input:    []bool{true, true, true, true, false, false, false, false, true, true, true, true, false, false, false, false},
			expected: []byte{0x0F, 0x0F},
		},
		{
			name:     "Unaligned - requires padding",
			input:    []bool{true, true, true},
			expected: []byte{0x07},
		},
		{
			name:     "Single bit",
			input:    []bool{true},
			expected: []byte{0x01},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BoolSliceToByteSlice(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("BoolSliceToByteSlice() = %v, want %v", result, tt.expected)
			}
		})
	}
}
