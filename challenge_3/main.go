package main

import (
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

const (
	inputBuffer = "1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736"
)

var charFreq = map[string]float64{
	"a": 0.0651738, "b": 0.0124248, "c": 0.0217339, "d": 0.0349835, "e": 0.1041442, "f": 0.0197881, "g": 0.0158610,
	"h": 0.0492888, "i": 0.0558094, "j": 0.0009033, "k": 0.0050529, "l": 0.0331490, "m": 0.0202124, "n": 0.0564513,
	"o": 0.0596302, "p": 0.0137645, "q": 0.0008606, "r": 0.0497563, "s": 0.0515760, "t": 0.0729357, "u": 0.0225134,
	"v": 0.0082903, "w": 0.0171272, "x": 0.0013692, "y": 0.0145984, "z": 0.0007836, " ": 0.1918182,
}

func main() {
	cipher, _ := hex.DecodeString(inputBuffer)

	bruteForce(cipher)
	automated(cipher)
}

// Print out all options, sort through visually to find
func bruteForce(src []byte) {
	for i := 0; i < 256; i++ {
		key := byte(i)

		dst := make([]byte, len(src))

		for k := 0; k < len(src); k++ {
			dst[k] = src[k] ^ key
		}

		fmt.Println(string(dst))
		fmt.Printf("\n\n")

	}
}

type Result struct {
	character     int
	score         float64
	decryptedText string
}

func (r *Result) decrypt(cipher []byte) {
	dst := make([]byte, len(cipher))
	for i := 0; i < len(cipher); i++ {
		dst[i] = cipher[i] ^ byte(r.character)
	}

	r.decryptedText = string(dst)
}

func (r *Result) calculateScore() {
	var score float64 = 0

	for i := 0; i < len(r.decryptedText); i++ {
		char := strings.ToLower(string(r.decryptedText[i]))
		val, ok := charFreq[char]

		if !ok {
			continue
		}

		score += val
	}

	r.score = score
}

func newResult(c int) Result {
	return Result{
		character:     c,
		score:         0,
		decryptedText: "",
	}
}

func automated(src []byte) {
	results := make([]Result, 256)

	for i := 0; i < 256; i++ {
		r := newResult(i)
		r.decrypt(src)
		r.calculateScore()
		results[i] = r
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})

	for i := 0; i < len(results); i++ {
		fmt.Println(string(results[i].character), results[i].score, ":", results[i].decryptedText)
	}
}
