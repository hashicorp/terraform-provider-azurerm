package network

import "github.com/hashicorp/go-cty/cty"

func ShouldSuppressVnetAddressSpaceDiffForTest(rawConfig cty.Value) bool {
	return shouldSuppressVnetAddressSpaceDiff(rawConfig)
}
