package testhelpers

import (
	"math/rand"
	"time"
)

func RandomInt() int {
	reseed()
	return rand.New(rand.NewSource(time.Now().UnixNano())).Int()
}

func RandomString() string {
	size := 5
	charSet := "abcdefghijklmnopqrstuvwxyz0123456789"

	reseed()
	result := make([]byte, size)
	for i := 0; i < size; i++ {
		result[i] = charSet[rand.Intn(len(charSet))]
	}
	return string(result)
}

func reseed() {
	rand.Seed(time.Now().UTC().UnixNano())
}
