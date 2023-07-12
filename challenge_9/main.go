package main

import (
	"encoding/hex"
	"fmt"
)

const (
	input = "YELLOW SUBMARINE"

	input2 = "test"
)

func main() {
	blocksize := 20

	src := []byte(input)
	fmt.Println(hex.EncodeToString(src))

	dst := pad(src, blocksize)

	fmt.Println(string(dst), len(string(dst)))
}

func pad(src []byte, n int) []byte {
	dst := make([]byte, n)
	np := n - len(src)

	padding := make([]byte, np)
	for i := range padding {
		padding[i] = byte(np)
	}

	dst = append(src, padding...)

	return dst
}
