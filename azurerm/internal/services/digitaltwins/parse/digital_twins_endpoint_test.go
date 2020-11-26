package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = DigitalTwinsEndpointId{}

func TestDigitalTwinsEndpointIDFormatter(t *testing.T) {
	actual := NewDigitalTwinsEndpointID("12345678-1234-9876-4563-123456789012", "group1", "instance1", "endpoint1").ID("")
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.DigitalTwins/digitalTwinsInstances/instance1/endpoints/endpoint1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestDigitalTwinsEndpointID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *DigitalTwinsEndpointId
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
			// missing DigitalTwinsInstanceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.DigitalTwins/",
			Error: true,
		},

		{
			// missing value for DigitalTwinsInstanceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.DigitalTwins/digitalTwinsInstances/",
			Error: true,
		},

		{
			// missing EndpointName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.DigitalTwins/digitalTwinsInstances/instance1/",
			Error: true,
		},

		{
			// missing value for EndpointName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.DigitalTwins/digitalTwinsInstances/instance1/endpoints/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.DigitalTwins/digitalTwinsInstances/instance1/endpoints/endpoint1",
			Expected: &DigitalTwinsEndpointId{
				SubscriptionId:           "12345678-1234-9876-4563-123456789012",
				ResourceGroup:            "group1",
				DigitalTwinsInstanceName: "instance1",
				EndpointName:             "endpoint1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/GROUP1/PROVIDERS/MICROSOFT.DIGITALTWINS/DIGITALTWINSINSTANCES/INSTANCE1/ENDPOINTS/ENDPOINT1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := DigitalTwinsEndpointID(v.Input)
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
		if actual.DigitalTwinsInstanceName != v.Expected.DigitalTwinsInstanceName {
			t.Fatalf("Expected %q but got %q for DigitalTwinsInstanceName", v.Expected.DigitalTwinsInstanceName, actual.DigitalTwinsInstanceName)
		}
		if actual.EndpointName != v.Expected.EndpointName {
			t.Fatalf("Expected %q but got %q for EndpointName", v.Expected.EndpointName, actual.EndpointName)
		}
	}
}
