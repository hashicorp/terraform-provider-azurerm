// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azure_test

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

func TestIsRetired(t *testing.T) {
	testData := []struct {
		year     int
		month    time.Month
		day      int
		expected bool
	}{
		{
			year:     2000,
			month:    time.December,
			day:      31,
			expected: true,
		},
		{
			year:     3000,
			month:    time.December,
			day:      31,
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing: %s", time.Date(v.year, v.month, v.day, 0, 0, 0, 0, time.UTC))

		actual, _ := azure.IsRetired(v.year, v.month, v.day)
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
