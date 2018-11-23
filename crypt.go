package main

func decrypt(b, key []byte) []byte {
	if len(b) != 8 {
		return b
	}
	phase1 := shuffle(b)
	phase2 := xor(phase1, key)
	phase3 := shift(phase2)
	ctmp := offset()
	result := calc(phase3, ctmp)
	return result
}

func calc(b, ctmp []byte) []byte {
	res := make([]byte, 8)
	for i := range b {
		res[i] = (0xFF + b[i] - ctmp[i] + 0x01) & 0xFF
	}
	return res
}

func offset() []byte {
	offset := []byte{0x48, 0x74, 0x65, 0x6D, 0x70, 0x39, 0x39, 0x65} //"Htemp99e"
	res := make([]byte, 8)
	for i := range offset {
		res[i] = (offset[i]>>4 | offset[i]<<4) & 0xFF
	}
	return res
}

func shift(b []byte) []byte {
	res := make([]byte, 8)
	for i := range b {
		res[i] = (b[i]>>3 | b[(i-1+8)%8]<<5) & 0xFF
	}
	return res
}

func xor(b, key []byte) []byte {
	res := make([]byte, 8)
	for i := range b {
		res[i] = b[i] ^ key[i]
	}
	return res
}

func shuffle(b []byte) []byte {
	assignNum := []int{2, 4, 0, 7, 1, 6, 5, 3}
	res := make([]byte, 8)
	for i, v := range assignNum {
		res[i] = b[v]
	}
	return res
}
