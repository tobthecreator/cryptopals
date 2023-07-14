package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	mrand "math/rand"
)

const (
	BLOCKSIZE = 16
)

func main() {
	dst := encryptUnknown([]byte("hello"))

	fmt.Println(dst)
}

func generateRandKey() ([]byte, error) {
	key := make([]byte, BLOCKSIZE)

	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func genPad(n int) []byte {
	padding := make([]byte, n)
	for i := range padding {
		padding[i] = byte(n)
	}

	return padding
}

func coinFlip() bool {
	b := make([]byte, 1)
	_, err := rand.Read(b)
	if err != nil {
		return false
	}

	// Check LSB to see if the num is odd or even
	return b[0]&1 == 1
}

func generateRandomPadding() (dst []byte) {
	min := 5
	max := 10

	// Max number is the range, then offset by the min to bring it into range
	n := mrand.Intn(max-min) + min
	p := make([]byte, n)

	rand.Read(p)

	return p
}

func encryptUnknown(src []byte) []byte {
	plaintext := make([]byte, len(src)) // arrays are pbr, not pbv. avoid side effects
	copy(plaintext, src)

	randKey, _ := generateRandKey()
	useCbc := coinFlip()

	basePadding := genPad(((len(src)/BLOCKSIZE + 1) * BLOCKSIZE) - len(src))

	plaintext = append(plaintext, basePadding...)

	c, _ := aes.NewCipher(randKey)

	ciphertext := make([]byte, len(plaintext))

	if useCbc {
		randIV, _ := generateRandKey()
		e := cipher.NewCBCEncrypter(c, randIV)

		e.CryptBlocks(ciphertext, plaintext)
	} else {

		for i := 0; i < len(plaintext)/BLOCKSIZE; i += BLOCKSIZE {
			c.Encrypt(ciphertext[i:i+BLOCKSIZE], plaintext[i:i+BLOCKSIZE])
		}
	}

	leftPadding := generateRandomPadding()
	ciphertext = append(leftPadding, ciphertext...)

	rightPadding := generateRandomPadding()
	ciphertext = append(ciphertext, rightPadding...)

	return ciphertext
}
