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
	h, _ := hex.DecodeString(plaintext)
	b := make([]byte, hex.DecodedLen(len(h)))
	hex.Decode(b, h)

	return Cipher{
		Plaintext: plaintext,
		Hex:       h,
		Binary:    b,
	}
}

func (c *Cipher) Print() {
	fmt.Println(c.Hex)
}
