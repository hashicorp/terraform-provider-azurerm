// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package point_to_site_vpn_gateway_test

import (
	"math/rand"
)

func randString() string {
	charSet := "abcdefghijklmnopqrstuvwxyz"
	strlen := 5
	result := make([]byte, strlen)
	for i := range strlen {
		result[i] = charSet[rand.Intn(len(charSet))]
	}
	return string(result)
}
