// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func TestStorageDiscoveryScopeReplacementPath(t *testing.T) {
	base := storageDiscoveryScopeTestData("TestScope", []interface{}{"Microsoft.Storage/storageAccounts"}, []interface{}{}, map[string]interface{}{})
	additional := storageDiscoveryScopeTestData("AdditionalScope", []interface{}{"Microsoft.Storage/storageAccounts"}, []interface{}{}, map[string]interface{}{})

	testCases := []struct {
		name     string
		old      []interface{}
		new      []interface{}
		expected string
	}{
		{
			name:     "add scope",
			old:      []interface{}{base},
			new:      []interface{}{base, additional},
			expected: "",
		},
		{
			name:     "remove scope",
			old:      []interface{}{base, additional},
			new:      []interface{}{base},
			expected: "",
		},
		{
			name: "change resource types",
			old:  []interface{}{base},
			new: []interface{}{
				storageDiscoveryScopeTestData("TestScope", []interface{}{"Microsoft.Storage/storageAccounts", "Microsoft.Storage/storageAccounts/blobServices"}, []interface{}{}, map[string]interface{}{}),
			},
			expected: "scope.0.resource_types",
		},
		{
			name: "change tag keys",
			old:  []interface{}{base},
			new: []interface{}{
				storageDiscoveryScopeTestData("TestScope", []interface{}{"Microsoft.Storage/storageAccounts"}, []interface{}{"environment"}, map[string]interface{}{}),
			},
			expected: "scope.0.tag_keys_only",
		},
		{
			name: "change tags",
			old:  []interface{}{base},
			new: []interface{}{
				storageDiscoveryScopeTestData("TestScope", []interface{}{"Microsoft.Storage/storageAccounts"}, []interface{}{}, map[string]interface{}{"environment": "test"}),
			},
			expected: "scope.0.tags",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := storageDiscoveryScopeReplacementPath(testCase.old, testCase.new)
			if actual != testCase.expected {
				t.Fatalf("expected %q but got %q", testCase.expected, actual)
			}
		})
	}
}

func TestStorageDiscoveryWorkspaceScopeDiffRequiresReplacement(t *testing.T) {
	base := storageDiscoveryScopeConfig("TestScope", []interface{}{"Microsoft.Storage/storageAccounts"}, []interface{}{}, map[string]interface{}{})
	additional := storageDiscoveryScopeConfig("AdditionalScope", []interface{}{"Microsoft.Storage/storageAccounts"}, []interface{}{}, map[string]interface{}{})
	filtered := storageDiscoveryScopeConfig("FilteredScope", []interface{}{"Microsoft.Storage/storageAccounts"}, []interface{}{"environment"}, map[string]interface{}{"tier": "production"})
	unfiltered := storageDiscoveryScopeConfig("FilteredScope", []interface{}{"Microsoft.Storage/storageAccounts"}, []interface{}{}, map[string]interface{}{})

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
			name:        "insert scope before an existing scope",
			oldScopes:   []interface{}{filtered},
			newScopes:   []interface{}{base, filtered},
			requiresNew: false,
		},
		{
			name:        "remove scope before an existing scope",
			oldScopes:   []interface{}{base, filtered},
			newScopes:   []interface{}{filtered},
			requiresNew: false,
		},
		{
			name:        "reorder scopes",
			oldScopes:   []interface{}{base, filtered},
			newScopes:   []interface{}{filtered, base},
			requiresNew: false,
		},
		{
			name:        "remove scope and change shifted scope filters",
			oldScopes:   []interface{}{base, filtered},
			newScopes:   []interface{}{unfiltered},
			requiresNew: true,
		},
		{
			name:        "insert scope and change shifted scope filters",
			oldScopes:   []interface{}{filtered, base},
			newScopes:   []interface{}{additional, unfiltered, base},
			requiresNew: true,
		},
		{
			name:        "reorder and change scope filters",
			oldScopes:   []interface{}{filtered, base},
			newScopes:   []interface{}{base, unfiltered},
			requiresNew: true,
		},
		{
			name:        "remove scope and change shifted scope tags",
			oldScopes:   []interface{}{base, storageDiscoveryScopeConfig("AdditionalScope", []interface{}{"Microsoft.Storage/storageAccounts"}, []interface{}{}, map[string]interface{}{"tier": "production"})},
			newScopes:   []interface{}{additional},
			requiresNew: true,
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

func TestStorageDiscoveryWorkspaceRootsDiff(t *testing.T) {
	const subscriptionID = "/subscriptions/aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	const otherSubscriptionID = "/subscriptions/bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"

	testCases := []struct {
		name      string
		roots     []interface{}
		wantError bool
	}{
		{
			name:  "subscription",
			roots: []interface{}{subscriptionID},
		},
		{
			name:  "resource groups",
			roots: []interface{}{subscriptionID + "/resourceGroups/first", subscriptionID + "/resourceGroups/second"},
		},
		{
			name:  "subscription and unrelated resource group",
			roots: []interface{}{subscriptionID, otherSubscriptionID + "/resourceGroups/test"},
		},
		{
			name:      "subscription and child resource group",
			roots:     []interface{}{subscriptionID, subscriptionID + "/resourceGroups/test"},
			wantError: true,
		},
		{
			name:      "subscription and child with different subscription casing",
			roots:     []interface{}{"/subscriptions/AAAAAAAA-AAAA-AAAA-AAAA-AAAAAAAAAAAA", subscriptionID + "/resourceGroups/test"},
			wantError: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			wrapper := sdk.NewResourceWrapper(StorageDiscoveryWorkspaceResource{})
			resource, err := wrapper.Resource()
			if err != nil {
				t.Fatalf("building resource: %v", err)
			}

			config := storageDiscoveryWorkspaceConfig([]interface{}{
				storageDiscoveryScopeConfig("TestScope", []interface{}{"Microsoft.Storage/storageAccounts"}, []interface{}{}, map[string]interface{}{}),
			})
			config["workspace_roots"] = testCase.roots
			_, err = resource.Diff(context.Background(), nil, terraform.NewResourceConfigRaw(config), &clients.Client{})
			if testCase.wantError {
				if err == nil || !strings.Contains(err.Error(), "cannot specify both subscription ID") {
					t.Fatalf("expected a workspace root overlap error, got: %v", err)
				}
			} else if err != nil {
				t.Fatalf("building diff: %v", err)
			}
		})
	}
}

func TestStorageDiscoveryWorkspaceRootValidation(t *testing.T) {
	validator := StorageDiscoveryWorkspaceResource{}.Arguments()["workspace_roots"].Elem.(*pluginsdk.Schema).ValidateFunc
	testCases := []struct {
		name  string
		value string
		valid bool
	}{
		{
			name:  "subscription",
			value: "/subscriptions/00000000-0000-0000-0000-000000000000",
			valid: true,
		},
		{
			name:  "resource group",
			value: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test",
			valid: true,
		},
		{
			name:  "storage account",
			value: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test/providers/Microsoft.Storage/storageAccounts/test",
		},
		{
			name:  "invalid ID",
			value: "not-an-id",
		},
		{
			name: "empty ID",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, errors := validator(testCase.value, "workspace_roots")
			if (len(errors) == 0) != testCase.valid {
				t.Fatalf("expected valid %t, got errors: %v", testCase.valid, errors)
			}
		})
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
