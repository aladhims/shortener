package base

import (
	"strings"
)

// ALPHABET is all possible characters
const ALPHABET = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

// BASE is the base for conversion based on the length of ALPHABET
const BASE = uint64(len(ALPHABET))

func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

// Encode will return back the converted version of the origin with base64 (the length of the ALPHABET)
func Encode(origin uint64) (encoded string) {
	if origin == 0 {
		encoded = string(ALPHABET[0])
		return
	}

	for origin > 0 {
		index := origin % BASE
		encoded += string(ALPHABET[index])
		origin /= BASE
	}

	encoded = reverse(encoded)
	return
}

// Decode will return back the original version of the encoded string with base64 (the length of the ALPHABET)
func Decode(encoded string) (origin uint64) {
	for _, char := range encoded {
		index := strings.Index(ALPHABET, string(char))
		origin = (BASE * origin) + uint64(index)
	}
	return
}
