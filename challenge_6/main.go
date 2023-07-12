package main

import (
	"bufio"
	"cryptopals/cipher"
	"cryptopals/decrypt"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

const (
	input = `I'm back and I'm ringin' the bell
A rockin' on the mike while the fly girls yell
In ecstasy in the back of me
Well that's my DJ Deshay cuttin' all them Z's
Hittin' hard and the girlies goin' crazy
Vanilla's on the mike, man I'm not lazy.

I'm lettin' my drug kick in
It controls my mouth and I begin
To just let it flow, let my concepts go
My posse's to the side yellin', Go Vanilla Go!

Smooth 'cause that's the way I will be
And if you don't give a damn, then
Why you starin' at me
So get off 'cause I control the stage
There's no dissin' allowed
I'm in my own phase
The girlies sa y they love me and that is ok
And I can dance better than any kid n' play

Stage 2 -- Yea the one ya' wanna listen to
It's off my head so let the beat play through
So I can funk it up and make it sound good
1-2-3 Yo -- Knock on some wood
For good luck, I like my rhymes atrocious
Supercalafragilisticexpialidocious
I'm an effect and that you can bet
I can take a fly girl and make her wet.

I'm like Samson -- Samson to Delilah
There's no denyin', You can try to hang
But you'll keep tryin' to get my style
Over and over, practice makes perfect
But not if you're a loafer.

You'll get nowhere, no place, no time, no girls
Soon -- Oh my God, homebody, you probably eat
Spaghetti with a spoon! Come on and say it!

VIP. Vanilla Ice yep, yep, I'm comin' hard like a rhino
Intoxicating so you stagger like a wino
So punks stop trying and girl stop cryin'
Vanilla Ice is sellin' and you people are buyin'
'Cause why the freaks are jockin' like Crazy Glue
Movin' and groovin' trying to sing along
All through the ghetto groovin' this here song
Now you're amazed by the VIP posse.

Steppin' so hard like a German Nazi
Startled by the bases hittin' ground
There's no trippin' on mine, I'm just gettin' down
Sparkamatic, I'm hangin' tight like a fanatic
You trapped me once and I thought that
You might have it
So step down and lend me your ear
'89 in my time! You, '90 is my year.

You're weakenin' fast, YO! and I can tell it
Your body's gettin' hot, so, so I can smell it
So don't be mad and don't be sad
'Cause the lyrics belong to ICE, You can call me Dad
You're pitchin' a fit, so step back and endure
Let the witch doctor, Ice, do the dance to cure
So come up close and don't be square
You wanna battle me -- Anytime, anywhere

You thought that I was weak, Boy, you're dead wrong
So come on, everybody and sing this song

Say -- Play that funky music Say, go white boy, go white boy go
play that funky music Go white boy, go white boy, go
Lay down and boogie and play that funky music till you die.

Play that funky music Come on, Come on, let me hear
Play that funky music white boy you say it, say it
Play that funky music A little louder now
Play that funky music, white boy Come on, Come on, Come on
Play that funky music`
)

func main() {
	b := readFile2()

	results := make([]decrypt.ApproximatorResult, 39)
	for i := 2; i <= 40; i++ {
		results[i-2] = decrypt.ApproximatorResult{
			Keysize:                i,
			AverageHammingDistance: 0,
		}
	}

	fmt.Println("len(b)", len(b))
	// For each keysize
	for i, r := range results {
		distances := []float64{}

		if r.Keysize == 0 {
			fmt.Println("hey boss it's zero")
			continue
		}

		// Calculate the Hamming Distance for each subset of pairs of size r.Keysize
		for j := 0; j < len(b)-(len(b)%r.Keysize)-r.Keysize; j += r.Keysize {
			// for j := 0; j < 29*3; j += r.Keysize {

			// First block of size r.Keysize, block n
			a := b[j : j+r.Keysize]

			// Second following block of size r.Keysize, block n+1
			b := b[j+r.Keysize : j+r.Keysize*2]

			normalizedDistance := float64(decrypt.HammingDistance(a, b)) / float64(r.Keysize)
			distances = append(distances, normalizedDistance)
		}

		fmt.Println(r.Keysize, "\n\n", distances, "\n\n")
		// Normalize distances by keysize, then add to running sum
		var sum float64 = 0
		for j := 0; j < len(distances); j++ {
			sum += distances[j]
		}

		// Calculate the average
		results[i].AverageHammingDistance = sum / float64(len(distances))

	}

	// Sort results so that the shortest Hamming Distance is first
	sort.Slice(results, func(i, j int) bool {
		return results[i].AverageHammingDistance < results[j].AverageHammingDistance
	})

	for _, r := range results {
		fmt.Println("key: ", r.Keysize, "d: ", r.AverageHammingDistance)
	}

	keysize := results[0].Keysize

	// Crack it
	// So I need to break everything up into sections of size keysize
	fullIterations := len(b) / keysize
	fmt.Println("keysize: ", keysize)
	fmt.Println("fi: ", fullIterations)
	singleBitCiphers := make([][]byte, keysize)
	for i := range singleBitCiphers {
		singleBitCiphers[i] = []byte{}
	}

	// full iterations
	for i := 0; i < fullIterations; i++ {
		cipherIndex := i % keysize
		// cipherLen := len(singleBitCiphers[cipherIndex])

		singleBitCiphers[cipherIndex] = append(singleBitCiphers[cipherIndex], b[i])
	}

	decrypters := make([]decrypt.SingleBitXorDecrypter, keysize)
	bestResults := make([]decrypt.Result, keysize)
	for i, sbc := range singleBitCiphers {
		fmt.Println(sbc)

		decrypters[i] = decrypt.NewSingleBitXorDecrypter(
			cipher.NewCipher(
				hex.EncodeToString(sbc),
			),
		)

		decrypters[i].Decrypt()

		bestResults[i] = decrypters[i].Results[0]
	}

	for _, result := range bestResults {
		fmt.Println(result.Key, result.DecryptedText)
	}

}

// func main2() {
// 	cipher := readFile2()
// 	// cipher := readFile()
// 	// fmt.Println(cipher.Plaintext)
// 	ka := decrypt.NewKeysizeApproximator(cipher)
// 	ka.Evaluate()
// 	for _, r := range ka.Results {
// 		fmt.Println("key: ", r.Keysize, "d: ", r.AverageHammingDistance)
// 	}

// 	key := []byte{0x54, 0x65, 0x72, 0x6d, 0x69, 0x6e, 0x61, 0x74, 0x6f, 0x72, 0x20, 0x58, 0x3a, 0x20, 0x42, 0x72, 0x69, 0x6e, 0x67, 0x20, 0x74, 0x68, 0x65, 0x20, 0x6e, 0x6f, 0x69, 0x73, 0x65}
// 	key2 := []byte{84, 101, 114, 109, 105, 110, 97, 116, 111, 114, 32, 88, 58, 32, 66, 114, 105, 110, 103, 32, 116, 104, 101, 32, 110, 111, 105, 115, 101}
// 	fmt.Println(key)
// 	fmt.Println(string(key))
// 	fmt.Println(string(key2))

// 	e := encrypt.NewRepeatingKeyXorEncrypter(input, string(key))
// 	fmt.Println(e.Cipher.Plaintext[0:50])
// 	fmt.Println(cipher.Plaintext[0:50])
// 	fmt.Println(e.Cipher.Hex[0:50])
// 	fmt.Println(cipher.Hex[0:50])
// 	fmt.Println(e.Cipher.Binary[0:50])
// 	fmt.Println(cipher.Binary[0:50])
// 	fmt.Println(e.Cipher.Plaintext[:66] == cipher.Plaintext[:66])

// 	fmt.Println(len(cipher.Plaintext) == len(e.Cipher.Plaintext), len(cipher.Plaintext), len(e.Cipher.Plaintext))

// 	for i := 0; i < len(e.Cipher.Binary); i++ {
// 		j := i % len(key)
// 		fmt.Print(string(e.Cipher.Binary[i] ^ key[j]))
// 	}

// 	for i := 0; i < len(cipher.Binary); i++ {
// 		j := i % len(key)
// 		fmt.Print(string(cipher.Binary[i] ^ key[j]))
// 	}

// 	// fmt.Println(string(cipher.Plaintext[0:70]))

// 	// fmt.Println(string(e.Cipher.Plaintext[0:70]))
// }

func readFile() cipher.Cipher {
	file, _ := os.Open("./challenge_6/data.txt")
	defer file.Close()

	var lines []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// fmt.Println(line, "\n")
		lines = append(lines, line)
	}

	src := strings.Join(lines, "")

	// fmt.Println(src)

	b, _ := base64.StdEncoding.DecodeString(src)

	// fmt.Println(b)
	// fmt.Println(string(b))

	// // Decode the base64-encoded hex string
	// fmt.Println(string(data))
	// h, err := base64.StdEncoding.decode
	// if err != nil {
	// 	fmt.Println("Error decoding base64:", err)
	// }

	return cipher.NewCipher(string(b))
}

// Maybe my assumption is just bad. I've been expecting the decoded base64 to be
func readFile2() []byte {
	data, _ := ioutil.ReadFile("./challenge_6/data.txt")

	decodedBytes, _ := base64.StdEncoding.DecodeString(string(data))
	fmt.Println(decodedBytes)
	fmt.Println(string(decodedBytes))

	key := []byte{0x54, 0x65, 0x72, 0x6d, 0x69, 0x6e, 0x61, 0x74, 0x6f, 0x72, 0x20, 0x58, 0x3a, 0x20, 0x42, 0x72, 0x69, 0x6e, 0x67, 0x20, 0x74, 0x68, 0x65, 0x20, 0x6e, 0x6f, 0x69, 0x73, 0x65}

	for i := 0; i < len(decodedBytes); i++ {
		j := i % len(key)
		fmt.Print(string(decodedBytes[i] ^ key[j]))
	}

	c := cipher.NewCipher(string(decodedBytes))

	fmt.Println("\n\n-----------\n\n")

	for i := 0; i < len(c.Binary); i++ {
		j := i % len(key)
		fmt.Print(string(c.Binary[i] ^ key[j]))
	}

	// c2 := readFile()

	// fmt.Println("binary equal?", reflect.DeepEqual(c2.Binary, decodedBytes))
	// fmt.Println("hex equal?", reflect.DeepEqual(c2.Hex, c.Hex))
	// fmt.Println("plaintext equal?", reflect.DeepEqual(c2.Plaintext, c.Plaintext))
	// return c

	return decodedBytes
}
