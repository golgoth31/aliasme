package utils

import (
	"math/rand"
)

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	numberRunes = []rune("0123456789")
	rnd         = rand.New(rand.NewSource(rand.Int63()))
)

// GenerateRandomString generates a random string of specified length
func GenerateRandomString(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rnd.Intn(len(letterRunes))]
	}
	return string(b)
}

// GenerateRandomNumber generates a random numeric string of specified length
func GenerateRandomNumber(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = numberRunes[rnd.Intn(len(numberRunes))]
	}
	return string(b)
}
