// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-04-01/redisenterprise"
)

func TestManagedRedisCacheSkuName(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		expectedSku   *redisenterprise.Sku
		expectedError bool
		errorContains string
	}{
		// Valid cases - Enterprise SKUs with capacity
		{
			name:  "valid Enterprise E10 with capacity 2",
			input: "Enterprise_E10-2",
			expectedSku: &redisenterprise.Sku{
				Name:     redisenterprise.SkuName("Enterprise_E10"),
				Capacity: pointer.To(int64(2)),
			},
			expectedError: false,
		},

		// Valid cases - EnterpriseFlash SKUs with capacity
		{
			name:  "valid EnterpriseFlash F300 with capacity 3",
			input: "EnterpriseFlash_F300-3",
			expectedSku: &redisenterprise.Sku{
				Name:     redisenterprise.SkuName("EnterpriseFlash_F300"),
				Capacity: pointer.To(int64(3)),
			},
			expectedError: false,
		},

		// Valid cases - Non-Enterprise SKUs without capacity
		{
			name:  "valid Balanced B5 without capacity",
			input: "Balanced_B5",
			expectedSku: &redisenterprise.Sku{
				Name:     redisenterprise.SkuName("Balanced_B5"),
				Capacity: nil,
			},
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ManagedRedisCacheSkuName(tc.input)

			if tc.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tc.errorContains != "" && !strings.Contains(err.Error(), tc.errorContains) {
					t.Errorf("Expected error to contain %q, got %q", tc.errorContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}

				if result == nil {
					t.Errorf("Expected result but got nil")
					return
				}

				if tc.expectedSku != nil {
					if result.Name != tc.expectedSku.Name {
						t.Errorf("Expected Name %q, got %q", tc.expectedSku.Name, result.Name)
					}

					if tc.expectedSku.Capacity == nil {
						if result.Capacity != nil {
							t.Errorf("Expected Capacity to be nil, got %v", *result.Capacity)
						}
					} else {
						if result.Capacity == nil {
							t.Errorf("Expected Capacity %v, got nil", *tc.expectedSku.Capacity)
						} else if *result.Capacity != *tc.expectedSku.Capacity {
							t.Errorf("Expected Capacity %v, got %v", *tc.expectedSku.Capacity, *result.Capacity)
						}
					}
				}
			}
		})
	}
}
