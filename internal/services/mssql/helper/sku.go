package helper

import (
	"strings"
)

// CompareDatabaseSkuServiceTiers returns true if sku1 has a higher service tier than sku2
func CompareDatabaseSkuServiceTiers(sku1, sku2 string) bool {
	// These are intentionally short so that both forms can be matched, e.g. "S1" or "Standard"
	order := []string{
		"", "B", "S", "P",
	}

	var index1, index2 int
	for i, v := range order {
		if strings.HasPrefix(strings.ToLower(sku1), strings.ToLower(v)) {
			index1 = i
		}
		if strings.HasPrefix(strings.ToLower(sku2), strings.ToLower(v)) {
			index2 = i
		}
	}

	return index1 > 0 && index2 > 0 && index1 > index2
}
