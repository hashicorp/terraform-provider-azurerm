// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"testing"

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
