package main

import (
	"cryptopals/decrypt"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
)

var charFreq = map[string]float64{
	"a": 0.0651738, "b": 0.0124248, "c": 0.0217339, "d": 0.0349835, "e": 0.1041442, "f": 0.0197881, "g": 0.0158610,
	"h": 0.0492888, "i": 0.0558094, "j": 0.0009033, "k": 0.0050529, "l": 0.0331490, "m": 0.0202124, "n": 0.0564513,
	"o": 0.0596302, "p": 0.0137645, "q": 0.0008606, "r": 0.0497563, "s": 0.0515760, "t": 0.0729357, "u": 0.0225134,
	"v": 0.0082903, "w": 0.0171272, "x": 0.0013692, "y": 0.0145984, "z": 0.0007836, " ": 0.1918182,
}

func main() {
	src := readFile()

	keysize := approximateKeysize(src)

	// Now we go through and group together the bytes that would be worked on by each part of the key into their own ciphers
	ciphers := groupSingleBitXorCiphers(src, keysize)

	// Single-Bit Xor Decrypt each individual cipher, then add the best scored result to []key
	key := make([]byte, len(ciphers))
	for i, cipher := range ciphers {
		key[i] = decryptSingleBitXor(cipher)
	}

	decryptWithKeyAndPrint(src, key)
}

func approximateKeysize(src []byte) int {
	approxResults := make([]decrypt.ApproximatorResult, 39)
	for i := 2; i <= 40; i++ {
		approxResults[i-2] = decrypt.ApproximatorResult{
			Keysize:                i,
			AverageHammingDistance: 0,
		}
	}

	// For each keysize
	for i, r := range approxResults {
		distances := []float64{}

		if r.Keysize == 0 {
			continue
		}

		// Calculate the Hamming Distance for each subset of pairs of size r.Keysize
		for j := 0; j < len(src)-(len(src)%r.Keysize)-r.Keysize; j += r.Keysize {

			// First block of size r.Keysize, block n
			a := src[j : j+r.Keysize]

			// Second following block of size r.Keysize, block n+1
			b := src[j+r.Keysize : j+r.Keysize*2]

			normalizedDistance := float64(decrypt.HammingDistance(a, b)) / float64(r.Keysize)
			distances = append(distances, normalizedDistance)
		}

		// Normalize distances by keysize, then add to running sum
		var sum float64 = 0
		for j := 0; j < len(distances); j++ {
			sum += distances[j]
		}

		// Calculate the average
		approxResults[i].AverageHammingDistance = sum / float64(len(distances))

	}

	// Sort approxResults so that the shortest Hamming Distance is first
	sort.Slice(approxResults, func(i, j int) bool {
		return approxResults[i].AverageHammingDistance < approxResults[j].AverageHammingDistance
	})

	return approxResults[0].Keysize
}

func groupSingleBitXorCiphers(src []byte, keysize int) [][]byte {
	ciphers := make([][]byte, keysize)
	for i := range ciphers {
		ciphers[i] = []byte{}
	}

	for i := 0; i < len(src); i++ {
		cipherIndex := i % keysize
		ciphers[cipherIndex] = append(ciphers[cipherIndex], src[i])
	}

	return ciphers
}

func decryptSingleBitXor(src []byte) byte {
	l := len(src)

	type ScoringResult struct {
		Score float64
		Key   byte
	}

	scoringResults := make([]ScoringResult, 256)
	for i := 0; i < 256; i++ {

		dst := make([]byte, l)
		for j := 0; j < len(src); j++ {
			dst[j] = src[j] ^ byte(i)
		}

		str := string(dst)
		score := 0.0

		lowercase := strings.ToLower(str)
		for i := 0; i < len(str); i++ {
			char := string(lowercase[i])
			val, ok := charFreq[char]

			if !ok {
				continue
			}

			score += val
		}

		scoringResults[i] = ScoringResult{
			Key:   byte(i),
			Score: score,
		}
	}

	sort.Slice(scoringResults, func(i, j int) bool {
		return scoringResults[i].Score > scoringResults[j].Score
	})

	return scoringResults[0].Key
}

func decryptWithKeyAndPrint(b []byte, k []byte) {
	for i, v := range b {
		fmt.Print(string(v ^ k[i%len(k)]))
	}
}

// Maybe my assumption is just bad. I've been expecting the decoded base64 to be
func readFile() []byte {
	data, _ := ioutil.ReadFile("./challenge_6/data.txt")

	b, _ := base64.StdEncoding.DecodeString(string(data))

	return b
}
