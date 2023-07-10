// Convert hex to base64

// The string:
// 49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d

// Should produce:
// SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t

// Make that happen.

// Always operate on raw bytes, never on encoded strings. Only use hex and base64 for pretty-printing.

package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

const (
	hextable = "0123456789abcdef"

	/*
		The ASCII table is 256 unique characters that we can use to represent basic string text

		ASCII includes A-Z, a-z and 0-9.  a-z and 0-9 (case insensitive) are valid hex characters

		The reverse hex table functions by using taking in an byte value as an index for a position in the table.

		To look up the number 0, we'd look at index 48.

		If an index returned from the reverseHexTable is greater than x0f, than we know it is an invalid hex value.
	*/
	reverseHexTable = "" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\xff\xff\xff\xff\xff\xff" + // [0-9] are represented at 48-57
		"\xff\x0a\x0b\x0c\x0d\x0e\x0f\xff\xff\xff\xff\xff\xff\xff\xff\xff" + // [A-F] are represented at 65-70
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\x0a\x0b\x0c\x0d\x0e\x0f\xff\xff\xff\xff\xff\xff\xff\xff\xff" + // [a-f] are represented at 97-102
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff"
	hexString    = "49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d"
	base64String = "SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t"

	base64StandardEncoding = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
)

func main() {

	// Custom rolled version, closely following Golang source code
	fmt.Println("TOB: ", base64String == binaryToBase64(hexToBinary(hexString)))

	// Golang stdlib version
	hexBytes, _ := hex.DecodeString(hexString)
	base64.StdEncoding.EncodeToString(hexBytes)

	fmt.Println("Golang:", base64String == base64.StdEncoding.EncodeToString(hexBytes))
}

func hexToBinary(s string) []byte {
	// []byte(s) converts the input string into it's corresponding uint8 values.
	// If you print src, you'll see an array of uint8 numbers representing the decimal version of each byte.
	// We're going to use these decimals to look up hexidecimal values.
	src := []byte(s)
	dst := []byte(s)

	// i represents the location we're writing to in our output byte array
	// j is the location we're scanning in the hexadecimal
	i, j := 0, 1

	// A single hexadecimal character represents 4 bits, and there are 8 bits to a byte, our output.
	// Every 4 bits is a nibble.
	// So we're going to navigate the hex input string 2 nibbles at a time
	for ; j < len(src); j += 2 {
		msc := src[j-1] // most significant hex character
		lsc := src[j]   // least significant hex character

		a := reverseHexTable[msc] // most significant nibble
		b := reverseHexTable[lsc] // least significant nibble

		// verify these are acceptable single-char hex values
		if a > 0x0f || b > 0x0f {
			return dst[:i]
		}

		/*
			To combine two hex values into a byte, we need to bitshift the most significant nibble into position,
			then or it with the least significant nibble.

			If a = 0100 and b=0111, we need 0100 0111

			So we bit shift a, effectively padding it with four 0s.

			a -> 0100 0000

			Then we OR a with b, so 01000000 | 00000111

			We don't need to pad b's left side, it's understood to be all zeros

			The output of this is our first byte!
		*/
		dst[i] = (a << 4) | b
		i++ // Iterate to our next byte.
	}

	if len(src)%2 == 1 {
		if reverseHexTable[dst[j-1]] > 0x0f {
			return dst[:i]
		}
		return dst[:i]
	}

	return dst[:i]
}

func binaryToBase64(src []byte) string {
	const paddingChar rune = '='

	srcLength := len(src)
	base64Length := (srcLength + 2) / 3 * 4

	dst := make([]byte, base64Length)

	if srcLength == 0 {
		return ""
	}

	/*
		Base64 values can be represented by 6 bits.

		Our inputs are an array of bytes, each representing 8 bits.

		To convert cleanly from 8 bits to 6 bits, we need to find some intermediary representation.

		8 * 3 is 24.  6 * 4 is 24.

		So for every 3 bytes we take from our source, we can generate 4 base64 output values.

		Once we have our 3 bytes in a 24 bit pattern, we can easily create the 4 output base64 values with bitwise operations and masking.
	*/

	iDst, iSrc := 0, 0

	// We will loop over as many sequences of 3 as we can, then pad the rest after the for loop
	numTriplets := (srcLength / 3) * 3

	// For as long as we have sequences of 3
	for iSrc < numTriplets {
		// Convert 3x 8bit source bytes into 4 bytes
		val24bit := uint(src[iSrc])<<16 | uint(src[iSrc+1])<<8 | uint(src[iSrc+2])

		dst[iDst+0] = base64StandardEncoding[val24bit>>18&0x3F]
		dst[iDst+1] = base64StandardEncoding[val24bit>>12&0x3F]
		dst[iDst+2] = base64StandardEncoding[val24bit>>6&0x3F]
		dst[iDst+3] = base64StandardEncoding[val24bit&0x3F]

		iSrc += 3
		iDst += 4
	}

	// Add padding
	remain := srcLength - iSrc
	if remain != 0 {
		// Add the remaining small block
		val := uint(src[iSrc+0]) << 16
		if remain == 2 {
			val |= uint(src[iSrc+1]) << 8
		}

		dst[iDst+0] = base64StandardEncoding[val>>18&0x3F]
		dst[iDst+1] = base64StandardEncoding[val>>12&0x3F]

		switch remain {
		case 2:
			dst[iDst+2] = base64StandardEncoding[val>>6&0x3F]
			dst[iDst+3] = byte(paddingChar)
		case 1:
			dst[iDst+2] = byte(paddingChar)
			dst[iDst+3] = byte(paddingChar)
		}
	}

	// Convert our 64bit encoding into ASCII characters
	return string(dst)
}
