package utils

func ByteToBoolArray(b byte) [8]bool {
	var boolArray [8]bool
	for i := 0; i < 8; i++ {
		boolArray[i] = (b & (1 << i)) != 0
	}
	return boolArray
}

func bitsToByte(mei [8]bool) byte {
	var result uint8
	for _, b := range mei {
		result >>= 1
		if b {
			result |= 0b10000000
		}
	}
	return result
}

func BoolArrayToByteArray(boolArray []bool) []byte {
	var bools []bool
	// Expand the array to the next highest multiple of 8
	if len(boolArray)%8 != 0 {
		bools = make([]bool, len(boolArray)+8-len(boolArray)%8)
	} else {
		bools = make([]bool, len(boolArray))
	}
	for i, b := range boolArray {
		bools[i] = b
	}
	data := make([]byte, len(bools)/8)
	for i := 0; i < len(bools); i = i + 8 {
		byteVal := bitsToByte([8]bool(bools[i : i+8]))
		data[i/8] = byteVal
	}
	return data
}
