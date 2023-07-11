package decrypt

import (
	"cryptopals/cipher"
	"sort"
	"strings"
)

type SingleBitXorDecrypter struct {
	Cipher       cipher.Cipher
	HighestScore float64
	Keys         [][]byte
	Results      []Result
}

var charFreq = map[string]float64{
	"a": 0.0651738, "b": 0.0124248, "c": 0.0217339, "d": 0.0349835, "e": 0.1041442, "f": 0.0197881, "g": 0.0158610,
	"h": 0.0492888, "i": 0.0558094, "j": 0.0009033, "k": 0.0050529, "l": 0.0331490, "m": 0.0202124, "n": 0.0564513,
	"o": 0.0596302, "p": 0.0137645, "q": 0.0008606, "r": 0.0497563, "s": 0.0515760, "t": 0.0729357, "u": 0.0225134,
	"v": 0.0082903, "w": 0.0171272, "x": 0.0013692, "y": 0.0145984, "z": 0.0007836, " ": 0.1918182,
}

func NewSingleBitXorDecrypter(c cipher.Cipher) SingleBitXorDecrypter {
	keys := make([][]byte, 256)
	for i := range keys {
		keys[i] = make([]byte, len(c.Hex))
	}

	results := make([]Result, 256)

	d := SingleBitXorDecrypter{
		Cipher:       c,
		HighestScore: 0,
		Keys:         keys,
		Results:      results,
	}

	d.GenerateKeys()

	return d
}

func (d *SingleBitXorDecrypter) GenerateKeys() {
	keyLen := len(d.Cipher.Hex)

	for i := 0; i < 256; i++ {
		x := byte(i)
		k := make([]byte, keyLen)

		for j := 0; j < keyLen; j++ {
			k[j] = x
		}

		d.Keys[i] = k
	}
}

func (d *SingleBitXorDecrypter) Decrypt() {
	cipherLen := len(d.Cipher.Hex)

	for i := 0; i < len(d.Results); i++ {
		var result Result
		dst := make([]byte, cipherLen)

		for j := 0; j < cipherLen; j++ {
			dst[j] = d.Cipher.Hex[j] ^ d.Keys[i][j]
		}

		result.Key = d.Keys[i]
		result.DecryptedText = string(dst)
		result.Score = d.Score(result.DecryptedText)
		result.Plaintext = d.Cipher.Plaintext

		d.Results[i] = result
	}

	sort.Slice(d.Results, func(i, j int) bool {
		return d.Results[i].Score > d.Results[j].Score
	})

	d.HighestScore = d.Results[0].Score
}

func (d *SingleBitXorDecrypter) Score(s string) float64 {
	var score float64 = 0

	lowercase := strings.ToLower(s)

	for i := 0; i < len(s); i++ {
		char := string(lowercase[i])
		val, ok := charFreq[char]

		if !ok {
			continue
		}

		score += val
	}

	return score
}

// func (d *SingleBitXorDecrypter) PrintResults(n int) []Result {
// 	nResults := d.Results[:n]

// }
