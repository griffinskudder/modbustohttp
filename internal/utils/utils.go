package utils

func ByteToBoolArray(b byte) [8]bool {
	var boolArray [8]bool
	for i := 0; i < 8; i++ {
		boolArray[i] = (b & (1 << i)) != 0
	}
	return boolArray
}
