// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"testing"
)

func TestManagedRedisClusterSkuName(t *testing.T) {
	testCases := []struct {
		name        string
		input       interface{}
		expectError bool
	}{
		// Valid cases - Enterprise SKUs with capacity
		{
			name:        "valid Enterprise E10-2 SKU",
			input:       "Enterprise_E10-2",
			expectError: false,
		},

		// Valid cases - EnterpriseFlash SKUs with capacity
		{
			name:        "valid EnterpriseFlash F300-3 SKU",
			input:       "EnterpriseFlash_F300-3",
			expectError: false,
		},

		// Valid cases - Non-Enterprise SKUs without capacity
		{
			name:        "valid Balanced_B5 SKU",
			input:       "Balanced_B5",
			expectError: false,
		},

		// Invalid cases - Enterprise/EnterpriseFlash without capacity
		{
			name:        "invalid Enterprise E10 SKU, missing capacity",
			input:       "Enterprise_E10",
			expectError: true,
		},

		// Invalid cases - Non-Enterprise SKUs with capacity
		{
			name:        "invalid Balanced_B5-2 SKU, capacity not allowed",
			input:       "Balanced_B5-2",
			expectError: true,
		},

		// Invalid cases - Enterprise with odd capacity
		{
			name:        "invalid Enterprise E10-1 SKU, odd capacity",
			input:       "Enterprise_E10-1",
			expectError: true,
		},

		// Invalid cases - EnterpriseFlash with invalid capacity pattern
		{
			name:        "invalid EnterpriseFlash F300-4 SKU, invalid capacity pattern",
			input:       "EnterpriseFlash_F300-4",
			expectError: true,
		},
		{
			name:        "invalid EnterpriseFlash F700-6 SKU, invalid capacity pattern",
			input:       "EnterpriseFlash_F700-6",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, errors := ManagedRedisClusterSkuName(tc.input, "sku_name")

			if tc.expectError {
				if len(errors) == 0 {
					t.Errorf("Expected error but got none")
				}
			} else {
				if len(errors) > 0 {
					t.Errorf("Expected no errors but got: %v", errors)
				}
			}
		})
	}
}
