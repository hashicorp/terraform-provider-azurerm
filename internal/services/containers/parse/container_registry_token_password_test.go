// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = ContainerRegistryTokenPasswordId{}

func TestContainerRegistryTokenPasswordIDFormatter(t *testing.T) {
	actual := NewContainerRegistryTokenPasswordID("12345678-1234-9876-4563-123456789012", "resGroup1", "registry1", "token1", "password").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ContainerRegistry/registries/registry1/tokens/token1/passwords/password"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestContainerRegistryTokenPasswordID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *ContainerRegistryTokenPasswordId
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
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ContainerRegistry/",
			Error: true,
		},

		{
			// missing value for RegistryName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ContainerRegistry/registries/",
			Error: true,
		},

		{
			// missing TokenName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ContainerRegistry/registries/registry1/",
			Error: true,
		},

		{
			// missing value for TokenName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ContainerRegistry/registries/registry1/tokens/",
			Error: true,
		},

		{
			// missing PasswordName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ContainerRegistry/registries/registry1/tokens/token1/",
			Error: true,
		},

		{
			// missing value for PasswordName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ContainerRegistry/registries/registry1/tokens/token1/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ContainerRegistry/registries/registry1/tokens/token1/passwords/password",
			Expected: &ContainerRegistryTokenPasswordId{
				SubscriptionId: "12345678-1234-9876-4563-123456789012",
				ResourceGroup:  "resGroup1",
				RegistryName:   "registry1",
				TokenName:      "token1",
				PasswordName:   "password",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.CONTAINERREGISTRY/REGISTRIES/REGISTRY1/TOKENS/TOKEN1/PASSWORDS/PASSWORD",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ContainerRegistryTokenPasswordID(v.Input)
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
		if actual.TokenName != v.Expected.TokenName {
			t.Fatalf("Expected %q but got %q for TokenName", v.Expected.TokenName, actual.TokenName)
		}
		if actual.PasswordName != v.Expected.PasswordName {
			t.Fatalf("Expected %q but got %q for PasswordName", v.Expected.PasswordName, actual.PasswordName)
		}
	}
}
