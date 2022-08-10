package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = WebPubsubSharedPrivateLinkResourceId{}

func TestWebPubsubSharedPrivateLinkResourceIDFormatter(t *testing.T) {
	actual := NewWebPubsubSharedPrivateLinkResourceID("12345678-1234-9876-4563-123456789012", "resGroup1", "Webpubsub1", "resource1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.SignalRService/webPubSub/Webpubsub1/sharedPrivateLinkResources/resource1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestWebPubsubSharedPrivateLinkResourceID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *WebPubsubSharedPrivateLinkResourceId
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
			// missing WebPubSubName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.SignalRService/",
			Error: true,
		},

		{
			// missing value for WebPubSubName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.SignalRService/webPubSub/",
			Error: true,
		},

		{
			// missing SharedPrivateLinkResourceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.SignalRService/webPubSub/Webpubsub1/",
			Error: true,
		},

		{
			// missing value for SharedPrivateLinkResourceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.SignalRService/webPubSub/Webpubsub1/sharedPrivateLinkResources/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.SignalRService/webPubSub/Webpubsub1/sharedPrivateLinkResources/resource1",
			Expected: &WebPubsubSharedPrivateLinkResourceId{
				SubscriptionId:                "12345678-1234-9876-4563-123456789012",
				ResourceGroup:                 "resGroup1",
				WebPubSubName:                 "Webpubsub1",
				SharedPrivateLinkResourceName: "resource1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.SIGNALRSERVICE/WEBPUBSUB/WEBPUBSUB1/SHAREDPRIVATELINKRESOURCES/RESOURCE1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := WebPubsubSharedPrivateLinkResourceID(v.Input)
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
		if actual.WebPubSubName != v.Expected.WebPubSubName {
			t.Fatalf("Expected %q but got %q for WebPubSubName", v.Expected.WebPubSubName, actual.WebPubSubName)
		}
		if actual.SharedPrivateLinkResourceName != v.Expected.SharedPrivateLinkResourceName {
			t.Fatalf("Expected %q but got %q for SharedPrivateLinkResourceName", v.Expected.SharedPrivateLinkResourceName, actual.SharedPrivateLinkResourceName)
		}
	}
}
