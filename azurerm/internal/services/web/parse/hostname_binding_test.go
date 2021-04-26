package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = HostnameBindingId{}

func TestHostnameBindingIDFormatter(t *testing.T) {
	actual := NewHostnameBindingID("12345678-1234-9876-4563-123456789012", "mygroup1", "site1", "binding1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/mygroup1/providers/Microsoft.Web/sites/site1/hostNameBindings/binding1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestHostnameBindingID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *HostnameBindingId
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
			// missing SiteName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/mygroup1/providers/Microsoft.Web/",
			Error: true,
		},

		{
			// missing value for SiteName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/mygroup1/providers/Microsoft.Web/sites/",
			Error: true,
		},

		{
			// missing Name
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/mygroup1/providers/Microsoft.Web/sites/site1/",
			Error: true,
		},

		{
			// missing value for Name
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/mygroup1/providers/Microsoft.Web/sites/site1/hostNameBindings/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/mygroup1/providers/Microsoft.Web/sites/site1/hostNameBindings/binding1",
			Expected: &HostnameBindingId{
				SubscriptionId: "12345678-1234-9876-4563-123456789012",
				ResourceGroup:  "mygroup1",
				SiteName:       "site1",
				Name:           "binding1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/MYGROUP1/PROVIDERS/MICROSOFT.WEB/SITES/SITE1/HOSTNAMEBINDINGS/BINDING1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := HostnameBindingID(v.Input)
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
		if actual.SiteName != v.Expected.SiteName {
			t.Fatalf("Expected %q but got %q for SiteName", v.Expected.SiteName, actual.SiteName)
		}
		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}
	}
}
