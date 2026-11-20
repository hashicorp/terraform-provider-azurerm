// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package networkmanager

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func TestNetworkManagerScopeValidation(t *testing.T) {
	resource := sdk.WrappedResource(ManagerResource{})

	testCases := []struct {
		name  string
		scope map[string]any
		valid bool
	}{
		{
			name:  "missing scope IDs",
			scope: map[string]any{},
		},
		{
			name: "empty scope IDs",
			scope: map[string]any{
				"management_group_ids": []any{},
				"subscription_ids":     []any{},
			},
		},
		{
			name: "management group scope",
			scope: map[string]any{
				"management_group_ids": []any{"/providers/Microsoft.Management/managementGroups/test"},
			},
			valid: true,
		},
		{
			name: "subscription scope",
			scope: map[string]any{
				"subscription_ids": []any{"/subscriptions/00000000-0000-0000-0000-000000000000"},
			},
			valid: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			config := map[string]any{
				"name":                "test",
				"resource_group_name": "test",
				"location":            "westus2",
				"scope":               []any{testCase.scope},
			}

			diagnostics := resource.Validate(terraform.NewResourceConfigRaw(config))
			if diagnostics.HasError() == testCase.valid {
				t.Fatalf("expected valid %t, got diagnostics: %v", testCase.valid, diagnostics)
			}
		})
	}
}
