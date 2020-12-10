package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = AaaaRecordId{}

func TestAaaaRecordIDFormatter(t *testing.T) {
	actual := NewAaaaRecordID("12345678-1234-9876-4563-123456789012", "resGroup1", "privateDnsZone1", "eheh1").ID("")
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/privateDnsZones/privateDnsZone1/AAAA/eheh1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestAaaaRecordID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *AaaaRecordId
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
			// missing PrivateDnsZoneName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/",
			Error: true,
		},

		{
			// missing value for PrivateDnsZoneName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/privateDnsZones/",
			Error: true,
		},

		{
			// missing AAAAName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/privateDnsZones/privateDnsZone1/",
			Error: true,
		},

		{
			// missing value for AAAAName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/privateDnsZones/privateDnsZone1/AAAA/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/privateDnsZones/privateDnsZone1/AAAA/eheh1",
			Expected: &AaaaRecordId{
				SubscriptionId:     "12345678-1234-9876-4563-123456789012",
				ResourceGroup:      "resGroup1",
				PrivateDnsZoneName: "privateDnsZone1",
				AAAAName:           "eheh1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.NETWORK/PRIVATEDNSZONES/PRIVATEDNSZONE1/AAAA/EHEH1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := AaaaRecordID(v.Input)
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
		if actual.PrivateDnsZoneName != v.Expected.PrivateDnsZoneName {
			t.Fatalf("Expected %q but got %q for PrivateDnsZoneName", v.Expected.PrivateDnsZoneName, actual.PrivateDnsZoneName)
		}
		if actual.AAAAName != v.Expected.AAAAName {
			t.Fatalf("Expected %q but got %q for AAAAName", v.Expected.AAAAName, actual.AAAAName)
		}
	}
}
