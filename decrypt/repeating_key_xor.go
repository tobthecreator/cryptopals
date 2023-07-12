package decrypt

import (
	"cryptopals/cipher"
	"fmt"
	"sort"
)

type KeysizeApproximator struct {
	Cipher  cipher.Cipher
	Results []ApproximatorResult
}

type ApproximatorResult struct {
	Keysize                int
	AverageHammingDistance float64
}

func NewKeysizeApproximator(c cipher.Cipher) KeysizeApproximator {

	results := make([]ApproximatorResult, 39)
	for i := 2; i <= 40; i++ {
		results[i-2] = ApproximatorResult{
			Keysize:                i,
			AverageHammingDistance: 0,
		}
	}

	return KeysizeApproximator{
		Cipher:  c,
		Results: results,
	}
}

func (ka *KeysizeApproximator) Evaluate() {

	// For each keysize
	for i, r := range ka.Results {
		distances := []int{}

		if r.Keysize == 0 {
			fmt.Println("hey boss it's zero")
			continue
		}

		// Calculate the Hamming Distance for each subset of pairs of size r.Keysize
		for j := 0; j < len(ka.Cipher.Binary)-(len(ka.Cipher.Binary)%r.Keysize)-r.Keysize; j += r.Keysize {

			// First block of size r.Keysize, block n
			a := ka.Cipher.Binary[j : j+r.Keysize]

			// Second following block of size r.Keysize, block n+1
			b := ka.Cipher.Binary[j+r.Keysize : j+r.Keysize*2]

			distances = append(distances, HammingDistance(a, b))
		}

		// Normalize distances by keysize, then add to running sum
		var sum float64 = 0
		for j := 0; j < len(distances); j++ {
			sum += (float64(distances[j]) / float64(r.Keysize))
		}

		// Calculate the average
		ka.Results[i].AverageHammingDistance = sum / float64(len(distances))

	}

	// Sort results so that the shortest Hamming Distance is first
	sort.Slice(ka.Results, func(i, j int) bool {
		return ka.Results[i].AverageHammingDistance < ka.Results[j].AverageHammingDistance
	})
}

type RepeatingKeyXorDecrypter struct {
	Cipher       cipher.Cipher
	Approximator KeysizeApproximator
	// Results []Result
}

func (d *RepeatingKeyXorDecrypter) Decrypt() {

}
