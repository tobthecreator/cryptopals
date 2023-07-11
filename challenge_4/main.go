package main

import (
	"bufio"
	"cryptopals/cipher"
	"cryptopals/decrypt"
	"fmt"
	"os"
	"sort"
)

var charFreq = map[string]float64{
	"a": 0.0651738, "b": 0.0124248, "c": 0.0217339, "d": 0.0349835, "e": 0.1041442, "f": 0.0197881, "g": 0.0158610,
	"h": 0.0492888, "i": 0.0558094, "j": 0.0009033, "k": 0.0050529, "l": 0.0331490, "m": 0.0202124, "n": 0.0564513,
	"o": 0.0596302, "p": 0.0137645, "q": 0.0008606, "r": 0.0497563, "s": 0.0515760, "t": 0.0729357, "u": 0.0225134,
	"v": 0.0082903, "w": 0.0171272, "x": 0.0013692, "y": 0.0145984, "z": 0.0007836, " ": 0.1918182,
}

func main() {

	ciphers := readFile()
	fmt.Println(len(ciphers))
	decrypters := make([]decrypt.SingleBitXorDecrypter, len(ciphers))

	for i, cipher := range ciphers {
		decrypters[i] = decrypt.NewSingleBitXorDecrypter(cipher)
		decrypters[i].Decrypt()
	}

	resultSampleSize := 3

	// Take top X results of from each decrypter and aggregate them
	topResults := make([]decrypt.Result, (len(decrypters) * resultSampleSize))
	for i, d := range decrypters {
		tr := d.Results[:resultSampleSize]

		j := i * resultSampleSize

		for k := 0; k < len(tr); k++ {
			topResults[j+k] = tr[k]
		}
	}

	// Reorder them by the highest score
	sort.Slice(topResults, func(i, j int) bool {
		return topResults[i].Score > topResults[j].Score
	})

	// Sample the top X results
	for i := 0; i < resultSampleSize; i++ {
		r := topResults[i]

		fmt.Println(
			r.Key, r.Plaintext, r.Score, r.DecryptedText,
		)
	}
}

func readFile() []cipher.Cipher {
	file, _ := os.Open("./challenge_4/data.txt")
	defer file.Close()

	var lines []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	fmt.Println(len(lines))
	ciphers := make([]cipher.Cipher, len(lines))
	fmt.Println(len(ciphers))
	for i, l := range lines {
		c := cipher.NewCipher(l)
		ciphers[i] = c
	}

	return ciphers
}
