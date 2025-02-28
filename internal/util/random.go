package util

import (
	"math/rand"
	"strings"
)

func RandomStringDefaultCharset(minLength int, maxLength int) string {
	alphabet := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890#+-,.:;><&%$§!()")
	length := minLength
	if maxLength > minLength {
		length += rand.Intn(maxLength - minLength)
	}
	return randomString(length, alphabet)
}

func RandomStringAlphaNumerical(minLength, maxLength int) string {
	alphabet := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890")
	length := minLength
	if maxLength > minLength {
		length += rand.Intn(maxLength - minLength)
	}
	return randomString(length, alphabet)
}

func randomString(length int, alphabet []rune) string {
	alphabetSize := len(alphabet)
	var sb strings.Builder

	for range length {
		ch := alphabet[rand.Intn(alphabetSize)]
		sb.WriteRune(ch)
	}

	s := sb.String()
	return s
}
