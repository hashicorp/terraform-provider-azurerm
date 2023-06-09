package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = ResourceTagsId{}

func TestResourceTagsIDFormatter(t *testing.T) {
	actual := NewResourceTagsID("12345678-1234-9876-4563-123456789012", "test-rg", "Microsoft.Compute", "virtualMachines", "test-vm").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/test-rg/providers/Microsoft.Compute/virtualMachines/test-vm/providers/Microsoft.Resources/tags/default"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestResourceTagsID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *ResourceTagsId
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
			// missing Provider
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/test-rg/providers/",
			Error: true,
		},

		{
			// missing value for Namespace
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/test-rg/providers/Microsoft.Compute/",
			Error: true,
		},

		{
			// missing ResourceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/test-rg/providers/Microsoft.Compute/virtualMachines/",
			Error: true,
		},

		{
			// missing value for TagResource
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/test-rg/providers/Microsoft.Compute/virtualMachines/test-vm",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/test-rg/providers/Microsoft.Compute/virtualMachines/test-vm/providers/Microsoft.Resources/tags/default",
			Expected: &ResourceTagsId{
				SubscriptionId: "12345678-1234-9876-4563-123456789012",
				ResourceGroup:  "test-rg",
				Provider:       "Microsoft.Compute",
				Namespace:      "virtualMachines",
				ResourceName:   "test-vm",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/test-rg/PROVIDERS/MICROSOFT.COMPUTE/VIRTUALMACHINES/test-vm/PROVIDERS/MICROSOFT.RESOURCES/TAGS/DEFAULT",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ResourceTagsID(v.Input)
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
		if actual.Provider != v.Expected.Provider {
			t.Fatalf("Expected %q but got %q for Provider", v.Expected.Provider, actual.Provider)
		}
		if actual.Namespace != v.Expected.Namespace {
			t.Fatalf("Expected %q but got %q for Namespace", v.Expected.Namespace, actual.Namespace)
		}
		if actual.ResourceName != v.Expected.ResourceName {
			t.Fatalf("Expected %q but got %q for ResourceName", v.Expected.ResourceName, actual.ResourceName)
		}
	}
}
