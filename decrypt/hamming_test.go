package decrypt

import "testing"

func TestHammingDistance(t *testing.T) {
	t.Run("Provided test", func(t *testing.T) {
		s1 := "this is a test"
		s2 := "wokka wokka!!!"
		expected := 37

		result := HammingDistance(s1, s2)
		if result != expected {
			t.Errorf("Expected Hamming distance to be %d, but got %d", expected, result)
		}
	})

}
