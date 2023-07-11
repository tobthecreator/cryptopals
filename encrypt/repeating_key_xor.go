package encrypt

import (
	"cryptopals/cipher"
	"encoding/hex"
)

type RepeatingKeyXorEncrypter struct {
	Input  string
	Key    string
	Cipher cipher.Cipher
}

func NewRepeatingKeyXorEncrypter(i string, k string) RepeatingKeyXorEncrypter {
	e := RepeatingKeyXorEncrypter{
		Input: i,
		Key:   k,
	}

	e.Encrypt()

	return e
}

func (e *RepeatingKeyXorEncrypter) Encrypt() {
	src := []byte(e.Input) // binary version of src
	srcLen := len(src)

	keys := []byte(e.Key) // binary version of keys
	keysLen := len(keys)

	dst := make([]byte, srcLen)

	for i := 0; i < srcLen; i++ {
		ki := i % keysLen

		dst[i] = src[i] ^ keys[ki]
	}

	h := make([]byte, hex.EncodedLen(srcLen))
	hex.Encode(h, dst)
	e.Cipher = cipher.Cipher{
		Plaintext: hex.EncodeToString(dst),
		Hex:       h,
		Binary:    dst,
	}
}
