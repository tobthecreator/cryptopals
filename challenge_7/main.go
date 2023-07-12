package main

import (
	"crypto/aes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
)

func main() {
	src := readFile()
	aesKey := []byte("YELLOW SUBMARINE")

	c, err := aes.NewCipher(aesKey)

	if err != nil {
		fmt.Println("error", err)
		return
	}

	if len(src)%len(aesKey) != 0 {
		fmt.Println("Uneven key")
		return
	}

	dst := make([]byte, len(src))
	for i := 0; i < len(src)/len(aesKey); i += len(aesKey) {
		c.Decrypt(dst[i:i+len(aesKey)], src[i:i+len(aesKey)])
	}

	fmt.Println(string(dst))
}

func readFile() []byte {
	data, _ := ioutil.ReadFile("./challenge_7/data.txt")

	b, _ := base64.StdEncoding.DecodeString(string(data))

	return b
}
