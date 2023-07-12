package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"sort"
)

/*
	https://crypto.stackexchange.com/questions/20941/why-shouldnt-i-use-ecb-encryption/20946#20946

	detect whether two ECB-encrypted messages are identical;
	detect whether two ECB-encrypted messages share a common prefix;
	detect whether two ECB-encrypted messages share other common substrings, as long as those substrings are aligned at block boundaries; or
	detect whether (and where) a single ECB-encrypted message contains repetitive data (such as long runs of spaces or null bytes, repeated header fields or coincidentally repeated phrases in text).

	I have a hunch that I can solve this with only one of these criteria, but I'm going to write something for all of them I think

*/

const (
	KEYSIZE = 16
)

func main() {
	src, lineLen, numLines := readFile()

	lineIndex := detectIdenticalBlocks(src, lineLen, numLines)

	fmt.Println(lineIndex) // 132 matches answers I'm finding online
}

func detectIdenticalBlocks(src []byte, lineLen int, numLines int) (lineIndex int) {
	type IdenticalBlocksScore struct {
		Score     int
		LineIndex int
	}
	scores := make([]IdenticalBlocksScore, numLines)

	// For every line, count the number of repeating blocks
	for i := 0; i < numLines; i++ {
		li := i * lineLen
		line := src[li : li+lineLen]

		repeatingBlocks := 0

		/*
			To count repeating blocks, we'll iterate through each block
			then compare it against all following blocks.

			We don't need to check previous blocks because those previous blocks will
			have searched for equality

			We'll count all equality collisions as a simple scoring system.  Properly
			encrypted files don't have this pattern, so the detection of any score > 0
			is pretty damning
		*/
		for j := 0; j < lineLen/KEYSIZE-1; j++ {
			blockIndex := j * KEYSIZE
			block := line[blockIndex : blockIndex+KEYSIZE]

			remainingBlocks := line[blockIndex+KEYSIZE:]
			for k := 0; k < len(remainingBlocks)/KEYSIZE; k++ {
				comparisonBlockIndex := k * KEYSIZE
				comparisonBlock := remainingBlocks[comparisonBlockIndex : comparisonBlockIndex+KEYSIZE]

				if reflect.DeepEqual(block, comparisonBlock) {
					repeatingBlocks++
				}
			}
		}

		scores[i] = IdenticalBlocksScore{
			Score:     repeatingBlocks,
			LineIndex: i,
		}
	}

	// Sort to put highest scores first
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Score > scores[j].Score
	})

	return scores[0].LineIndex
}

func readFile() (src []byte, lineLen int, numLines int) {
	filename := "./challenge_8/data.txt"

	// Length
	file, _ := os.Open(filename)
	ll := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ll = len(scanner.Text())
	}
	file.Close()

	// Bytes
	b, _ := ioutil.ReadFile(filename)

	// NumLines
	nl := len(b) / ll

	return b, ll, nl
}
