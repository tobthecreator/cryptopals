package main

import (
	"encoding/hex"
	"fmt"
)

const (
	bufferOne = "1c0111001f010100061a024b53535009181c"
	bufferTwo = "686974207468652062756c6c277320657965"
	output    = "746865206b696420646f6e277420706c6179"
)

func main() {
	h1, _ := hex.DecodeString(bufferOne)
	h2, _ := hex.DecodeString(bufferTwo)

	if len(h1) != len(h2) {
		return
	}

	dst := make([]byte, len(h1))
	fmt.Println(h1)
	fmt.Println(h2)

	for i := 0; i < len(h1); i++ {
		a := h1[i]
		b := h2[i]

		dst[i] = a ^ b
	}

	fmt.Println("[]bytes: ", dst)
	fmt.Println("hex: ", hex.EncodeToString(dst))

	fmt.Println("Output: ", hex.EncodeToString(dst) == output)

}
