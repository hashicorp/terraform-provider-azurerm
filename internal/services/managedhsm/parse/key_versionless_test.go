// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = KeyVersionlessId{}

func TestKeyVersionlessIDFormatter(t *testing.T) {
	actual := NewKeyVersionlessID("12345678-1234-9876-4563-123456789012", "resGroup1", "mhsm1", "key1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.KeyVault/managedHSMs/mhsm1/keys/key1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestKeyVersionlessID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *KeyVersionlessId
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
			// missing ManagedHSMName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.KeyVault/",
			Error: true,
		},

		{
			// missing value for ManagedHSMName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.KeyVault/managedHSMs/",
			Error: true,
		},

		{
			// missing KeyName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.KeyVault/managedHSMs/mhsm1/",
			Error: true,
		},

		{
			// missing value for KeyName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.KeyVault/managedHSMs/mhsm1/keys/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.KeyVault/managedHSMs/mhsm1/keys/key1",
			Expected: &KeyVersionlessId{
				SubscriptionId: "12345678-1234-9876-4563-123456789012",
				ResourceGroup:  "resGroup1",
				ManagedHSMName: "mhsm1",
				KeyName:        "key1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.KEYVAULT/MANAGEDHSMS/MHSM1/KEYS/KEY1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := KeyVersionlessID(v.Input)
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
		if actual.ManagedHSMName != v.Expected.ManagedHSMName {
			t.Fatalf("Expected %q but got %q for ManagedHSMName", v.Expected.ManagedHSMName, actual.ManagedHSMName)
		}
		if actual.KeyName != v.Expected.KeyName {
			t.Fatalf("Expected %q but got %q for KeyName", v.Expected.KeyName, actual.KeyName)
		}
	}
}
