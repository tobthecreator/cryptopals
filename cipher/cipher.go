package cipher

import (
	"encoding/hex"
	"fmt"
)

type Cipher struct {
	Plaintext string
	Hex       []byte
	Binary    []byte
}

func NewCipher(plaintext string) Cipher {
	h, _ := hex.DecodeString(plaintext) // my pain is coming from the fact that we expect plaintext to be a hex
	b := []byte(plaintext)

	return Cipher{
		Plaintext: plaintext,
		Hex:       h,
		Binary:    b,
	}
}

func (c *Cipher) Print() {
	fmt.Println(c.Hex)
}
