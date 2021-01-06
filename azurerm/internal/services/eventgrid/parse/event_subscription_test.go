package parse

import (
	"testing"
)

func TestEventGridEventSubscriptionId(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *EventSubscriptionId
	}{
		{
			Name:     "Empty",
			Input:    "",
			Expected: nil,
		},
		{
			Name:     "No Resource Groups Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/providers/Microsoft.EventGrid/eventSubscriptions/subscription1",
			Expected: nil,
		},
		{
			Name:  "Subscription Scope",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.EventGrid/eventSubscriptions/subscription1",
			Expected: &EventSubscriptionId{
				Name:  "subscription1",
				Scope: "/subscriptions/00000000-0000-0000-0000-000000000000",
			},
		},
		{
			Name:  "Resource Group",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.EventGrid/eventSubscriptions/subscription1",
			Expected: &EventSubscriptionId{
				Name:  "subscription1",
				Scope: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1",
			},
		},
		{
			Name:  "Storage Account Scope",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/storage1/providers/Microsoft.EventGrid/eventSubscriptions/subscription1",
			Expected: &EventSubscriptionId{
				Name:  "subscription1",
				Scope: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/storage1",
			},
		},
		{
			Name:  "Event Grid Domain Scope",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.EventGrid/domains/domain1/providers/Microsoft.EventGrid/eventSubscriptions/subscription1",
			Expected: &EventSubscriptionId{
				Name:  "subscription1",
				Scope: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.EventGrid/domains/domain1",
			},
		},
		{
			Name:  "Event Grid Topic Scope",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.EventGrid/topics/topic1/providers/Microsoft.EventGrid/eventSubscriptions/subscription1",
			Expected: &EventSubscriptionId{
				Name:  "subscription1",
				Scope: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.EventGrid/topics/topic1",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := EventSubscriptionID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}

		if actual.Scope != v.Expected.Scope {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expected.Scope, actual.Scope)
		}
	}
}
