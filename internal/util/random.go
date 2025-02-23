package util

import (
	"math/rand"
	"strings"
)

func randomStringDefaultCharset(minLength int, maxLength int) string {
	alphabet := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890#+-,.:;><&%$§!()")
	return randomString(minLength+rand.Intn(maxLength-minLength), alphabet)
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
