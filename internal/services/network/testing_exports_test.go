// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package network

import "github.com/hashicorp/go-cty/cty"

func ShouldSuppressVnetAddressSpaceDiffForTest(rawConfig cty.Value) bool {
	return shouldSuppressVnetAddressSpaceDiff(rawConfig)
}
