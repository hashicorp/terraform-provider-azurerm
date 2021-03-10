package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = LiveOutputId{}

func TestLiveOutputIDFormatter(t *testing.T) {
	actual := NewLiveOutputID("12345678-1234-9876-4563-123456789012", "resGroup1", "account1", "event1", "output1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Media/mediaservices/account1/liveevents/event1/liveoutputs/output1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestLiveOutputID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *LiveOutputId
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
			// missing MediaserviceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Media/",
			Error: true,
		},

		{
			// missing value for MediaserviceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Media/mediaservices/",
			Error: true,
		},

		{
			// missing LiveeventName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Media/mediaservices/account1/",
			Error: true,
		},

		{
			// missing value for LiveeventName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Media/mediaservices/account1/liveevents/",
			Error: true,
		},

		{
			// missing Name
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Media/mediaservices/account1/liveevents/event1/",
			Error: true,
		},

		{
			// missing value for Name
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Media/mediaservices/account1/liveevents/event1/liveoutputs/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Media/mediaservices/account1/liveevents/event1/liveoutputs/output1",
			Expected: &LiveOutputId{
				SubscriptionId:   "12345678-1234-9876-4563-123456789012",
				ResourceGroup:    "resGroup1",
				MediaserviceName: "account1",
				LiveeventName:    "event1",
				Name:             "output1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.MEDIA/MEDIASERVICES/ACCOUNT1/LIVEEVENTS/EVENT1/LIVEOUTPUTS/OUTPUT1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := LiveOutputID(v.Input)
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
		if actual.MediaserviceName != v.Expected.MediaserviceName {
			t.Fatalf("Expected %q but got %q for MediaserviceName", v.Expected.MediaserviceName, actual.MediaserviceName)
		}
		if actual.LiveeventName != v.Expected.LiveeventName {
			t.Fatalf("Expected %q but got %q for LiveeventName", v.Expected.LiveeventName, actual.LiveeventName)
		}
		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}
	}
}
