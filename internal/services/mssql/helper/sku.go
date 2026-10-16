// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package helper

import (
	"strconv"
	"strings"
)

// CompareDatabaseSkuScaleUp returns true if sku1 is a higher service tier or capacity than sku2.
func CompareDatabaseSkuScaleUp(sku1, sku2 string) bool {
	index1, capacity1 := databaseSkuTierAndCapacity(sku1)
	index2, capacity2 := databaseSkuTierAndCapacity(sku2)

	if index1 > 0 && index2 > 0 && index1 > index2 {
		return true
	}

	if index1 == 0 || index2 == 0 || index1 != index2 {
		return false
	}

	series1 := databaseSkuSeries(sku1)
	series2 := databaseSkuSeries(sku2)

	// NOTE: capacities are only compared within the same series, so capacity1 > capacity2 is
	// sufficient on its own. We intentionally do not require capacity2 > 0 here, otherwise a
	// scale up from a zero-capacity SKU such as `S0` would be missed (`databaseSkuCapacity("S0")` is 0).
	return series1 != "" && strings.EqualFold(series1, series2) && capacity1 > capacity2
}

func databaseSkuTierAndCapacity(sku string) (int, int) {
	// This order was observed to be enforced by the API. These are intentionally short so that
	// both forms can be matched for DTU tiers, e.g. "S1" or "Standard"
	order := []string{
		"", "B", "S", "GP", "P", "BC",
	}

	var index int
	for i, v := range order {
		if strings.HasPrefix(strings.ToLower(sku), strings.ToLower(v)) {
			index = i
		}
	}

	return index, databaseSkuCapacity(sku)
}

func databaseSkuSeries(sku string) string {
	if sku == "" {
		return ""
	}

	if i := strings.LastIndex(sku, "_"); i >= 0 {
		return sku[:i]
	}

	for i := len(sku) - 1; i >= 0; i-- {
		if sku[i] < '0' || sku[i] > '9' {
			return sku[:i+1]
		}
	}

	return sku
}

func databaseSkuCapacity(sku string) int {
	if sku == "" {
		return 0
	}

	var value string
	if i := strings.LastIndex(sku, "_"); i >= 0 {
		value = sku[i+1:]
	} else {
		for i := len(sku) - 1; i >= 0; i-- {
			if sku[i] < '0' || sku[i] > '9' {
				value = sku[i+1:]
				break
			}
		}
	}

	capacity, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}

	return capacity
}
