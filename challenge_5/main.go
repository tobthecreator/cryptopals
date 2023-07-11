package main

import (
	"cryptopals/encrypt"
	"fmt"
)

const (
	input = "Burning 'em, if you ain't quick and nimble\nI go crazy when I hear a cymbal"
)

func main() {
	e := encrypt.NewRepeatingKeyXorEncrypter(input, "ICE")

	fmt.Println(e.Cipher.Plaintext)
	fmt.Println(e.Cipher.Hex)
	fmt.Println(e.Cipher.Binary)
}
