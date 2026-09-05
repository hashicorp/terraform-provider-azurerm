// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func TestStorageDiscoveryScopesRequireReplacement(t *testing.T) {
	base := storageDiscoveryScopeTestData("TestScope", []interface{}{"Microsoft.Storage/storageAccounts"}, []interface{}{}, map[string]interface{}{})
	additional := storageDiscoveryScopeTestData("AdditionalScope", []interface{}{"Microsoft.Storage/storageAccounts"}, []interface{}{}, map[string]interface{}{})

	testCases := []struct {
		name     string
		old      []interface{}
		new      []interface{}
		expected bool
	}{
		{
			name:     "add scope",
			old:      []interface{}{base},
			new:      []interface{}{base, additional},
			expected: false,
		},
		{
			name:     "remove scope",
			old:      []interface{}{base, additional},
			new:      []interface{}{base},
			expected: false,
		},
		{
			name: "change resource types",
			old:  []interface{}{base},
			new: []interface{}{
				storageDiscoveryScopeTestData("TestScope", []interface{}{"Microsoft.Storage/storageAccounts", "Microsoft.Storage/storageAccounts/blobServices"}, []interface{}{}, map[string]interface{}{}),
			},
			expected: true,
		},
		{
			name: "change tag keys",
			old:  []interface{}{base},
			new: []interface{}{
				storageDiscoveryScopeTestData("TestScope", []interface{}{"Microsoft.Storage/storageAccounts"}, []interface{}{"environment"}, map[string]interface{}{}),
			},
			expected: true,
		},
		{
			name: "change tags",
			old:  []interface{}{base},
			new: []interface{}{
				storageDiscoveryScopeTestData("TestScope", []interface{}{"Microsoft.Storage/storageAccounts"}, []interface{}{}, map[string]interface{}{"environment": "test"}),
			},
			expected: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := storageDiscoveryScopesRequireReplacement(testCase.old, testCase.new)
			if actual != testCase.expected {
				t.Fatalf("expected %t but got %t", testCase.expected, actual)
			}
		})
	}
}

func TestStorageDiscoveryWorkspaceScopeDiffRequiresReplacement(t *testing.T) {
	base := storageDiscoveryScopeConfig("TestScope", []interface{}{"Microsoft.Storage/storageAccounts"}, []interface{}{}, map[string]interface{}{})
	additional := storageDiscoveryScopeConfig("AdditionalScope", []interface{}{"Microsoft.Storage/storageAccounts"}, []interface{}{}, map[string]interface{}{})

	testCases := []struct {
		name        string
		oldScopes   []interface{}
		newScopes   []interface{}
		requiresNew bool
	}{
		{
			name:        "add scope",
			oldScopes:   []interface{}{base},
			newScopes:   []interface{}{base, additional},
			requiresNew: false,
		},
		{
			name:        "remove scope",
			oldScopes:   []interface{}{base, additional},
			newScopes:   []interface{}{base},
			requiresNew: false,
		},
		{
			name:      "change resource types",
			oldScopes: []interface{}{base},
			newScopes: []interface{}{
				storageDiscoveryScopeConfig("TestScope", []interface{}{"Microsoft.Storage/storageAccounts", "Microsoft.Storage/storageAccounts/blobServices"}, []interface{}{}, map[string]interface{}{}),
			},
			requiresNew: true,
		},
		{
			name:      "change tag keys",
			oldScopes: []interface{}{base},
			newScopes: []interface{}{
				storageDiscoveryScopeConfig("TestScope", []interface{}{"Microsoft.Storage/storageAccounts"}, []interface{}{"environment"}, map[string]interface{}{}),
			},
			requiresNew: true,
		},
		{
			name:      "change tags",
			oldScopes: []interface{}{base},
			newScopes: []interface{}{
				storageDiscoveryScopeConfig("TestScope", []interface{}{"Microsoft.Storage/storageAccounts"}, []interface{}{}, map[string]interface{}{"environment": "test"}),
			},
			requiresNew: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			wrapper := sdk.NewResourceWrapper(StorageDiscoveryWorkspaceResource{})
			resource, err := wrapper.Resource()
			if err != nil {
				t.Fatalf("building resource: %v", err)
			}

			oldConfig := terraform.NewResourceConfigRaw(storageDiscoveryWorkspaceConfig(testCase.oldScopes))
			createDiff, err := resource.Diff(context.Background(), nil, oldConfig, &clients.Client{})
			if err != nil {
				t.Fatalf("building initial diff: %v", err)
			}

			attributes, err := createDiff.Apply(nil, resource.CoreConfigSchema())
			if err != nil {
				t.Fatalf("building initial state: %v", err)
			}

			state := &terraform.InstanceState{
				ID:         "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test/providers/Microsoft.StorageDiscovery/storageDiscoveryWorkspaces/test",
				Attributes: attributes,
			}
			newConfig := terraform.NewResourceConfigRaw(storageDiscoveryWorkspaceConfig(testCase.newScopes))
			diff, err := resource.Diff(context.Background(), state, newConfig, &clients.Client{})
			if err != nil {
				t.Fatalf("building update diff: %v", err)
			}

			if actual := diff.RequiresNew(); actual != testCase.requiresNew {
				t.Fatalf("expected replacement %t but got %t: %#v", testCase.requiresNew, actual, diff.Attributes)
			}
		})
	}
}

func TestValidateStorageDiscoveryScopes(t *testing.T) {
	first := storageDiscoveryScopeTestData("TestScope", []interface{}{"Microsoft.Storage/storageAccounts"}, []interface{}{}, map[string]interface{}{})
	second := storageDiscoveryScopeTestData("TestScope", []interface{}{"Microsoft.Storage/storageAccounts"}, []interface{}{"environment"}, map[string]interface{}{})
	unique := storageDiscoveryScopeTestData("AdditionalScope", []interface{}{"Microsoft.Storage/storageAccounts"}, []interface{}{}, map[string]interface{}{})

	if err := validateStorageDiscoveryScopes([]interface{}{first, unique}); err != nil {
		t.Fatalf("validating unique scope display names: %v", err)
	}

	if err := validateStorageDiscoveryScopes([]interface{}{first, second}); err == nil {
		t.Fatal("expected an error for duplicate scope display names")
	}
}

func storageDiscoveryScopeTestData(displayName string, resourceTypes, tagKeysOnly []interface{}, tags map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"display_name":   displayName,
		"resource_types": pluginsdk.NewSet(pluginsdk.HashString, resourceTypes),
		"tag_keys_only":  pluginsdk.NewSet(pluginsdk.HashString, tagKeysOnly),
		"tags":           tags,
	}
}

func storageDiscoveryWorkspaceConfig(scopes []interface{}) map[string]interface{} {
	return map[string]interface{}{
		"name":                "test",
		"resource_group_name": "test",
		"location":            "westus2",
		"workspace_roots":     []interface{}{`/subscriptions/00000000-0000-0000-0000-000000000000`},
		"scope":               scopes,
	}
}

func storageDiscoveryScopeConfig(displayName string, resourceTypes, tagKeysOnly []interface{}, tags map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"display_name":   displayName,
		"resource_types": resourceTypes,
		"tag_keys_only":  tagKeysOnly,
		"tags":           tags,
	}
}
