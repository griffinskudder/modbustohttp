package utils

// ByteToBoolSlice converts a byte to an array of 8 booleans representing each bit.
func ByteToBoolSlice(b byte) []bool {
	bools := make([]bool, 8)
	for i := 0; i < 8; i++ {
		bools[i] = (b & (1 << i)) != 0
	}
	return bools
}

// bitsToByte converts an array of 8 booleans to a byte.
func bitsToByte(bools [8]bool) byte {
	var result uint8
	for _, b := range bools {
		result >>= 1
		if b {
			result |= 0b10000000
		}
	}
	return result
}

// BoolSliceToByteSlice converts a slice of booleans to a slice of bytes.
// Each byte represents 8 booleans. If the length of the boolean array is not a multiple of 8,
// it will be padded with false values at the end.
func BoolSliceToByteSlice(boolsInput []bool) []byte {
	var output []byte
	if len(boolsInput)%8 != 0 {
		// Ensure the output slice has enough space to hold all bits, padding with false if necessary.
		output = make([]byte, (len(boolsInput)+8-len(boolsInput)%8)/8)
	} else {
		output = make([]byte, len(boolsInput)/8)
	}
	for i := 0; i < len(boolsInput); i = i + 8 {
		byteVal := bitsToByte([8]bool(boolsInput[i : i+8]))
		output[i/8] = byteVal
	}
	return output
}
