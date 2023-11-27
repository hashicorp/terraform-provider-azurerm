// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = ContainerRegistryCacheRuleId{}

func TestContainerRegistryCacheRuleIDFormatter(t *testing.T) {
	actual := NewContainerRegistryCacheRuleID("12345678-1234-9876-4563-123456789012", "group1", "registry1", "rule1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.ContainerRegistry/registries/registry1/cacheRules/rule1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestContainerRegistryCacheRuleID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *ContainerRegistryCacheRuleId
	}{

		{
			// empty
			Input: "",
			Error: true,
		},

		{
			// missing SubscriptionId
			Input: "/",
			Error: true,
		},

		{
			// missing value for SubscriptionId
			Input: "/subscriptions/",
			Error: true,
		},

		{
			// missing ResourceGroup
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/",
			Error: true,
		},

		{
			// missing value for ResourceGroup
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/",
			Error: true,
		},

		{
			// missing RegistryName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.ContainerRegistry/",
			Error: true,
		},

		{
			// missing value for RegistryName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.ContainerRegistry/registries/",
			Error: true,
		},

		{
			// missing CacheRuleName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.ContainerRegistry/registries/registry1/",
			Error: true,
		},

		{
			// missing value for CacheRuleName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.ContainerRegistry/registries/registry1/cacheRules/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.ContainerRegistry/registries/registry1/cacheRules/rule1",
			Expected: &ContainerRegistryCacheRuleId{
				SubscriptionId: "12345678-1234-9876-4563-123456789012",
				ResourceGroup:  "group1",
				RegistryName:   "registry1",
				CacheRuleName:  "rule1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/GROUP1/PROVIDERS/MICROSOFT.CONTAINERREGISTRY/REGISTRIES/REGISTRY1/CACHERULES/RULE1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ContainerRegistryCacheRuleID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %s", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

		if actual.SubscriptionId != v.Expected.SubscriptionId {
			t.Fatalf("Expected %q but got %q for SubscriptionId", v.Expected.SubscriptionId, actual.SubscriptionId)
		}
		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for ResourceGroup", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
		if actual.RegistryName != v.Expected.RegistryName {
			t.Fatalf("Expected %q but got %q for RegistryName", v.Expected.RegistryName, actual.RegistryName)
		}
		if actual.CacheRuleName != v.Expected.CacheRuleName {
			t.Fatalf("Expected %q but got %q for CacheRuleName", v.Expected.CacheRuleName, actual.CacheRuleName)
		}
	}
}
