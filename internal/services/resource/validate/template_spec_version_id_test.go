package validate

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import "testing"

func TestTemplateSpecVersionID(t *testing.T) {
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
			// missing TemplateSpecName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/templateSpecRG/providers/Microsoft.Resources/",
			Valid: false,
		},

		{
			// missing value for TemplateSpecName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/templateSpecRG/providers/Microsoft.Resources/templateSpecs/",
			Valid: false,
		},

		{
			// missing VersionName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/templateSpecRG/providers/Microsoft.Resources/templateSpecs/templateSpec1/",
			Valid: false,
		},

		{
			// missing value for VersionName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/templateSpecRG/providers/Microsoft.Resources/templateSpecs/templateSpec1/versions/",
			Valid: false,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/templateSpecRG/providers/Microsoft.Resources/templateSpecs/templateSpec1/versions/v1.0",
			Valid: true,
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/TEMPLATESPECRG/PROVIDERS/MICROSOFT.RESOURCES/TEMPLATESPECS/TEMPLATESPEC1/VERSIONS/V1.0",
			Valid: false,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := TemplateSpecVersionID(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
