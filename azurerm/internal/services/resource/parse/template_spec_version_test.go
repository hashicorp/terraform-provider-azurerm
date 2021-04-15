package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = TemplateSpecVersionId{}

func TestTemplateSpecVersionIDFormatter(t *testing.T) {
	actual := NewTemplateSpecVersionID("12345678-1234-9876-4563-123456789012", "templateSpecRG", "templateSpec1", "v1.0").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/templateSpecRG/providers/Microsoft.Resources/templateSpecs/templateSpec1/versions/v1.0"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestTemplateSpecVersionID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *TemplateSpecVersionId
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
			// missing TemplateSpecName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/templateSpecRG/providers/Microsoft.Resources/",
			Error: true,
		},

		{
			// missing value for TemplateSpecName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/templateSpecRG/providers/Microsoft.Resources/templateSpecs/",
			Error: true,
		},

		{
			// missing VersionName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/templateSpecRG/providers/Microsoft.Resources/templateSpecs/templateSpec1/",
			Error: true,
		},

		{
			// missing value for VersionName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/templateSpecRG/providers/Microsoft.Resources/templateSpecs/templateSpec1/versions/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/templateSpecRG/providers/Microsoft.Resources/templateSpecs/templateSpec1/versions/v1.0",
			Expected: &TemplateSpecVersionId{
				SubscriptionId:   "12345678-1234-9876-4563-123456789012",
				ResourceGroup:    "templateSpecRG",
				TemplateSpecName: "templateSpec1",
				VersionName:      "v1.0",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/TEMPLATESPECRG/PROVIDERS/MICROSOFT.RESOURCES/TEMPLATESPECS/TEMPLATESPEC1/VERSIONS/V1.0",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := TemplateSpecVersionID(v.Input)
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
		if actual.TemplateSpecName != v.Expected.TemplateSpecName {
			t.Fatalf("Expected %q but got %q for TemplateSpecName", v.Expected.TemplateSpecName, actual.TemplateSpecName)
		}
		if actual.VersionName != v.Expected.VersionName {
			t.Fatalf("Expected %q but got %q for VersionName", v.Expected.VersionName, actual.VersionName)
		}
	}
}
