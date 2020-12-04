package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = EventHubConsumerGroupId{}

func TestEventHubConsumerGroupIDFormatter(t *testing.T) {
	actual := NewEventHubConsumerGroupID("12345678-1234-9876-4563-123456789012", "group1", "namespace1", "eventhub1", "consumergroup1").ID("")
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.EventHub/namespaces/namespace1/eventhubs/eventhub1/consumergroups/consumergroup1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestEventHubConsumerGroupID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *EventHubConsumerGroupId
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
			// missing NamespaceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.EventHub/",
			Error: true,
		},

		{
			// missing value for NamespaceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.EventHub/namespaces/",
			Error: true,
		},

		{
			// missing EventhubName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.EventHub/namespaces/namespace1/",
			Error: true,
		},

		{
			// missing value for EventhubName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.EventHub/namespaces/namespace1/eventhubs/",
			Error: true,
		},

		{
			// missing ConsumergroupName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.EventHub/namespaces/namespace1/eventhubs/eventhub1/",
			Error: true,
		},

		{
			// missing value for ConsumergroupName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.EventHub/namespaces/namespace1/eventhubs/eventhub1/consumergroups/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.EventHub/namespaces/namespace1/eventhubs/eventhub1/consumergroups/consumergroup1",
			Expected: &EventHubConsumerGroupId{
				SubscriptionId:    "12345678-1234-9876-4563-123456789012",
				ResourceGroup:     "group1",
				NamespaceName:     "namespace1",
				EventhubName:      "eventhub1",
				ConsumergroupName: "consumergroup1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/GROUP1/PROVIDERS/MICROSOFT.EVENTHUB/NAMESPACES/NAMESPACE1/EVENTHUBS/EVENTHUB1/CONSUMERGROUPS/CONSUMERGROUP1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := EventHubConsumerGroupID(v.Input)
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
		if actual.NamespaceName != v.Expected.NamespaceName {
			t.Fatalf("Expected %q but got %q for NamespaceName", v.Expected.NamespaceName, actual.NamespaceName)
		}
		if actual.EventhubName != v.Expected.EventhubName {
			t.Fatalf("Expected %q but got %q for EventhubName", v.Expected.EventhubName, actual.EventhubName)
		}
		if actual.ConsumergroupName != v.Expected.ConsumergroupName {
			t.Fatalf("Expected %q but got %q for ConsumergroupName", v.Expected.ConsumergroupName, actual.ConsumergroupName)
		}
	}
}
