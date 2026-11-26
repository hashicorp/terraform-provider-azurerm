// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package arc_machine_extension_test

import (
	"math/rand"
)

const (
	LetterBytes  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	NumberBytes  = "1234567890"
	SpecialBytes = "!@#$%^()"
)

func generateRandomPassword(n int) string {
	b := make([]byte, n)
	for i := range b {
		r := rand.Int()
		switch r % 3 {
		case 0:
			b[i] = LetterBytes[rand.Intn(len(LetterBytes))]
		case 1:
			b[i] = SpecialBytes[rand.Intn(len(SpecialBytes))]
		case 2:
			b[i] = NumberBytes[rand.Intn(len(NumberBytes))]
		}
	}
	return string(b)
}
