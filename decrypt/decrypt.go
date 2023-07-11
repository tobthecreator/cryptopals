package decrypt

type Decrypter interface {
	Decrypt()
	Score()
	PrintResults(n int) []string
}

type Result struct {
	Key           []byte
	Score         float64
	DecryptedText string
	Plaintext     string
}
