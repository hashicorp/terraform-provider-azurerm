// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"testing"
)

func TestOrchestratedVirtualMachineScaleSetMigrateState(t *testing.T) {
	cases := map[string]struct {
		StateVersion       int
		InputAttributes    map[string]interface{}
		ExpectedVMSizes    []interface{}
		ExpectedVMSizesLen int
	}{
		"vm_sizes_migration_single_size": {
			StateVersion: 0,
			InputAttributes: map[string]interface{}{
				"id":                  "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testGroup/providers/Microsoft.Compute/virtualMachineScaleSets/testVMSS",
				"name":                "testVMSS",
				"resource_group_name": "testGroup",
				"location":            "West Europe",
				"sku_profile": []interface{}{
					map[string]interface{}{
						"allocation_strategy": "LowestPrice",
						"vm_sizes": []interface{}{
							"Standard_D2s_v3",
						},
					},
				},
			},
			ExpectedVMSizes: []interface{}{
				map[string]interface{}{
					"name": "Standard_D2s_v3",
				},
			},
			ExpectedVMSizesLen: 1,
		},
		"vm_sizes_migration_multiple_sizes": {
			StateVersion: 0,
			InputAttributes: map[string]interface{}{
				"id":                  "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testGroup/providers/Microsoft.Compute/virtualMachineScaleSets/testVMSS",
				"name":                "testVMSS",
				"resource_group_name": "testGroup",
				"location":            "West Europe",
				"sku_profile": []interface{}{
					map[string]interface{}{
						"allocation_strategy": "Prioritized",
						"vm_sizes": []interface{}{
							"Standard_D2s_v3",
							"Standard_D4s_v3",
							"Standard_D8s_v3",
						},
					},
				},
			},
			ExpectedVMSizes: []interface{}{
				map[string]interface{}{
					"name": "Standard_D2s_v3",
				},
				map[string]interface{}{
					"name": "Standard_D4s_v3",
				},
				map[string]interface{}{
					"name": "Standard_D8s_v3",
				},
			},
			ExpectedVMSizesLen: 3,
		},
		"no_sku_profile": {
			StateVersion: 0,
			InputAttributes: map[string]interface{}{
				"id":                  "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testGroup/providers/Microsoft.Compute/virtualMachineScaleSets/testVMSS",
				"name":                "testVMSS",
				"resource_group_name": "testGroup",
				"location":            "West Europe",
			},
			ExpectedVMSizes:    nil,
			ExpectedVMSizesLen: 0,
		},
		"empty_vm_sizes": {
			StateVersion: 0,
			InputAttributes: map[string]interface{}{
				"id":                  "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testGroup/providers/Microsoft.Compute/virtualMachineScaleSets/testVMSS",
				"name":                "testVMSS",
				"resource_group_name": "testGroup",
				"location":            "West Europe",
				"sku_profile": []interface{}{
					map[string]interface{}{
						"allocation_strategy": "LowestPrice",
						"vm_sizes":            []interface{}{},
					},
				},
			},
			ExpectedVMSizes:    []interface{}{},
			ExpectedVMSizesLen: 0,
		},
	}

	for testName, tc := range cases {
		t.Run(testName, func(t *testing.T) {
			migratedState, err := OrchestratedVirtualMachineScaleSetUpgradeV0ToV1{}.UpgradeFunc()(context.TODO(), tc.InputAttributes, nil)
			if err != nil {
				t.Fatalf("migration failed: %v", err)
			}

			// Check if sku_profile exists
			if tc.ExpectedVMSizesLen == 0 && tc.ExpectedVMSizes == nil {
				// No sku_profile case
				if _, exists := migratedState["sku_profile"]; !exists {
					return // Expected behavior for no sku_profile
				}
			}

			// Extract migrated vm_sizes
			skuProfileRaw, exists := migratedState["sku_profile"]
			if !exists {
				t.Fatalf("sku_profile not found in migrated state")
			}

			skuProfileList := skuProfileRaw.([]interface{})
			if len(skuProfileList) == 0 {
				t.Fatalf("sku_profile list is empty")
			}

			skuProfile := skuProfileList[0].(map[string]interface{})
			vmSizesRaw, exists := skuProfile["vm_sizes"]
			if !exists && tc.ExpectedVMSizesLen > 0 {
				t.Fatalf("vm_sizes not found in migrated sku_profile")
			}

			if tc.ExpectedVMSizesLen == 0 {
				vmSizesList := vmSizesRaw.([]interface{})
				if len(vmSizesList) != 0 {
					t.Fatalf("expected empty vm_sizes list, got %d items", len(vmSizesList))
				}
				return
			}

			vmSizesList := vmSizesRaw.([]interface{})

			// Check length
			if len(vmSizesList) != tc.ExpectedVMSizesLen {
				t.Fatalf("expected %d vm_sizes, got %d", tc.ExpectedVMSizesLen, len(vmSizesList))
			}

			// Check each migrated VM size
			for i, expectedVMSize := range tc.ExpectedVMSizes {
				expectedMap := expectedVMSize.(map[string]interface{})
				actualMap := vmSizesList[i].(map[string]interface{})

				if actualMap["name"] != expectedMap["name"] {
					t.Fatalf("vm_sizes[%d].name: expected %q, got %q", i, expectedMap["name"], actualMap["name"])
				}

				// Ensure rank field is not present in v0 migration (it's new in v1)
				if _, hasRank := actualMap["rank"]; hasRank {
					t.Fatalf("vm_sizes[%d] should not have rank field in v0->v1 migration", i)
				}
			}
		})
	}
}
