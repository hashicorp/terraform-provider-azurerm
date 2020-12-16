package parse

import (
	"testing"
)

func TestSystenTopicEventGridEventSubscriptionId(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *SystemTopicEventSubscriptionId
	}{
		{
			Name:     "Empty",
			Input:    "",
			Expected: nil,
		},
		{
			Name:     "No Resource Groups Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/providers/Microsoft.EventGrid/systemTopics/topic1/eventSubscriptions/subscription1",
			Expected: nil,
		},
		{
			Name:  "Event Grid System Topic Scope",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.EventGrid/systemTopics/topic1/eventSubscriptions/subscription1",
			Expected: &SystemTopicEventSubscriptionId{
				Name:          "subscription1",
				SystemTopic:   "topic1",
				ResourceGroup: "resGroup1",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := SystemTopicEventSubscriptionID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}

		if actual.SystemTopic != v.Expected.SystemTopic {
			t.Fatalf("Expected %q but got %q for System Topic", v.Expected.SystemTopic, actual.SystemTopic)
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
	}
}
