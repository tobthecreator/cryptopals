package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"io/ioutil"
)

const BLOCKSIZE = 16

func main() {
	// src := readFile()
	// dst := make([]byte, len(src))

	// fmt.Println(src)
	key := []byte("YELLOW SUBMARINE")
	iv := []byte("0000000000000000")
	plaintext := "I'm a fuckin' walkin' paradox. No, I'm not."

	padding := genPad(((len(plaintext)/BLOCKSIZE + 1) * BLOCKSIZE) - len(plaintext))
	b := append([]byte(plaintext), padding...)

	src := encrypt(b, key, iv)
	referenceDecrypt(src, key)
	decrypt(src, key, iv)

}

func breakIntoBlocks(src []byte) [][]byte {
	var blocks [][]byte

	for i := 0; i < len(src); i += BLOCKSIZE {
		end := i + BLOCKSIZE
		if end > len(src) {
			end = len(src)
		}

		blocks = append(blocks, src[i:end])
	}

	return blocks
}

func encrypt(src []byte, key []byte, iv []byte) []byte {
	c, _ := aes.NewCipher(key)

	plaintextBlocks := breakIntoBlocks(src)
	var cipherTextBlocks [][]byte
	var out []byte

	for i, plainTextBlock := range plaintextBlocks {
		xorBlock := make([]byte, BLOCKSIZE)

		if i == 0 {
			copy(xorBlock, iv)
			// use the iv
		} else {
			copy(xorBlock, cipherTextBlocks[i-1])
		}

		combinedBlock := xor(plainTextBlock, xorBlock)

		dst := make([]byte, len(plainTextBlock))
		c.Encrypt(dst, combinedBlock)

		// combined := xor(dst, xorBlock)
		cipherTextBlocks = append(cipherTextBlocks, dst)

		out = append(out, dst...)
	}

	return out
}

func decrypt(src, key, iv []byte) {
	c, _ := aes.NewCipher(key)

	cipherTextBlocks := breakIntoBlocks(src)
	var plainTextBlocks [][]byte
	var out []byte

	for i, cipherTextBlock := range cipherTextBlocks {
		dst := make([]byte, len(cipherTextBlock))
		c.Decrypt(dst, cipherTextBlock)

		xorBlock := make([]byte, BLOCKSIZE)

		if i == 0 {
			copy(xorBlock, iv)
			// use the iv
		} else {
			copy(xorBlock, cipherTextBlocks[i-1])
		}

		combined := xor(xorBlock, dst)

		plainTextBlocks = append(plainTextBlocks, combined)
		out = append(out, combined...)
	}

	fmt.Println(string(out))
}

func xor(b1 []byte, b2 []byte) []byte {
	dst := make([]byte, len(b1))

	for i := range dst {
		dst[i] = b1[i] ^ b2[i]
	}

	return dst
}

// Uses builtins
func referenceDecrypt(src []byte, key []byte) {
	c, _ := aes.NewCipher(key)
	d := cipher.NewCBCDecrypter(c, []byte("0000000000000000"))

	dst := make([]byte, len(src))

	d.CryptBlocks(dst, src)

	fmt.Println(string(dst))
}

func readFile() []byte {
	// stupid instructions didn't say it was b64 encoded :|
	data, _ := ioutil.ReadFile("./challenge_10/data.txt")
	b, _ := base64.StdEncoding.DecodeString(string(data))

	return b
}

func genPad(n int) []byte {
	padding := make([]byte, n)
	for i := range padding {
		padding[i] = byte(n)
	}

	return padding
}
