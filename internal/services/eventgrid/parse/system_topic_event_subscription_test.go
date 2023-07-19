// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = SystemTopicEventSubscriptionId{}

func TestSystemTopicEventSubscriptionIDFormatter(t *testing.T) {
	actual := NewSystemTopicEventSubscriptionID("12345678-1234-9876-4563-123456789012", "resGroup1", "systemTopic1", "subscription1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.EventGrid/systemTopics/systemTopic1/eventSubscriptions/subscription1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestSystemTopicEventSubscriptionID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *SystemTopicEventSubscriptionId
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
			// missing SystemTopicName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.EventGrid/",
			Error: true,
		},

		{
			// missing value for SystemTopicName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.EventGrid/systemTopics/",
			Error: true,
		},

		{
			// missing EventSubscriptionName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.EventGrid/systemTopics/systemTopic1/",
			Error: true,
		},

		{
			// missing value for EventSubscriptionName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.EventGrid/systemTopics/systemTopic1/eventSubscriptions/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.EventGrid/systemTopics/systemTopic1/eventSubscriptions/subscription1",
			Expected: &SystemTopicEventSubscriptionId{
				SubscriptionId:        "12345678-1234-9876-4563-123456789012",
				ResourceGroup:         "resGroup1",
				SystemTopicName:       "systemTopic1",
				EventSubscriptionName: "subscription1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.EVENTGRID/SYSTEMTOPICS/SYSTEMTOPIC1/EVENTSUBSCRIPTIONS/SUBSCRIPTION1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := SystemTopicEventSubscriptionID(v.Input)
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
		if actual.SystemTopicName != v.Expected.SystemTopicName {
			t.Fatalf("Expected %q but got %q for SystemTopicName", v.Expected.SystemTopicName, actual.SystemTopicName)
		}
		if actual.EventSubscriptionName != v.Expected.EventSubscriptionName {
			t.Fatalf("Expected %q but got %q for EventSubscriptionName", v.Expected.EventSubscriptionName, actual.EventSubscriptionName)
		}
	}
}
