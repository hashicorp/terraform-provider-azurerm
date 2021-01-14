package validate

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import "testing"

func TestActionGroupID(t *testing.T) {
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
			// missing Name
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/microsoft.insights/",
			Valid: false,
		},

		{
			// missing value for Name
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/microsoft.insights/actionGroups/",
			Valid: false,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/microsoft.insights/actionGroups/actionGroup1",
			Valid: true,
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/GROUP1/PROVIDERS/MICROSOFT.INSIGHTS/ACTIONGROUPS/ACTIONGROUP1",
			Valid: false,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := ActionGroupID(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
