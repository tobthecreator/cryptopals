package decrypt

import (
	"errors"
)

func HammingWeight(bytes []byte) int {
	weight := 0

	for _, b := range bytes {
		for b != 0 {
			// Check if the least significant bit is set
			if b&1 == 1 {
				weight++
			}
			// Right-shift to check the next bit
			b >>= 1
		}
	}

	return weight
}

func HammingDistance(a []byte, b []byte) int {
	getValue := func(arr []byte, index int) (byte, error) {
		if index < 0 || index >= len(arr) {
			return 0, errors.New("Index out of range")
		}

		return arr[index], nil
	}

	d := 0

	longest := len(a)

	if len(b) > len(a) {
		longest = len(b)
	}

	for i := 0; i < longest; i++ {
		va, aerr := getValue(a, i)
		vb, berr := getValue(b, i)

		if aerr != nil || berr != nil {
			continue
		}

		d += HammingWeight([]byte{va ^ vb})
	}

	return d
}
