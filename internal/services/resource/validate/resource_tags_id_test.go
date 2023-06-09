package validate

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import "testing"

func TestResourceTagsID(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{

		{
			// empty
			Input: "",
			Valid: false,
		},

		{
			// missing SubscriptionId
			Input: "/",
			Valid: false,
		},

		{
			// missing value for SubscriptionId
			Input: "/subscriptions/",
			Valid: false,
		},

		{
			// missing ResourceGroup
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/",
			Valid: false,
		},

		{
			// missing value for ResourceGroup
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/",
			Valid: false,
		},

		{
			// missing Provider
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/test-rg/providers/",
			Valid: false,
		},

		{
			// missing value for Namespace
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/test-rg/providers/Microsoft.Compute/",
			Valid: false,
		},

		{
			// missing ResourceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/test-rg/providers/Microsoft.Compute/virtualMachines/",
			Valid: false,
		},

		{
			// missing value for TagResource
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/test-rg/providers/Microsoft.Compute/virtualMachines/test-vm",
			Valid: false,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/test-rg/providers/Microsoft.Compute/virtualMachines/test-vm/providers/Microsoft.Resources/tags/default",
			Valid: true,
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/test-rg/PROVIDERS/MICROSOFT.COMPUTE/VIRTUALMACHINES/test-vm/PROVIDERS/MICROSOFT.RESOURCES/TAGS/DEFAULT",
			Valid: false,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := ResourceTagsID(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
