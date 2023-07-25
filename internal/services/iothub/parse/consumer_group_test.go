// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = ConsumerGroupId{}

func TestConsumerGroupIDFormatter(t *testing.T) {
	actual := NewConsumerGroupID("12345678-1234-9876-4563-123456789012", "resGroup1", "hub1", "events", "group1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Devices/iotHubs/hub1/eventHubEndpoints/events/consumerGroups/group1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestConsumerGroupID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *ConsumerGroupId
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
			// missing IotHubName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Devices/",
			Error: true,
		},

		{
			// missing value for IotHubName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Devices/iotHubs/",
			Error: true,
		},

		{
			// missing EventHubEndpointName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Devices/iotHubs/hub1/",
			Error: true,
		},

		{
			// missing value for EventHubEndpointName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Devices/iotHubs/hub1/eventHubEndpoints/",
			Error: true,
		},

		{
			// missing Name
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Devices/iotHubs/hub1/eventHubEndpoints/events/",
			Error: true,
		},

		{
			// missing value for Name
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Devices/iotHubs/hub1/eventHubEndpoints/events/consumerGroups/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Devices/iotHubs/hub1/eventHubEndpoints/events/consumerGroups/group1",
			Expected: &ConsumerGroupId{
				SubscriptionId:       "12345678-1234-9876-4563-123456789012",
				ResourceGroup:        "resGroup1",
				IotHubName:           "hub1",
				EventHubEndpointName: "events",
				Name:                 "group1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.DEVICES/IOTHUBS/HUB1/EVENTHUBENDPOINTS/EVENTS/CONSUMERGROUPS/GROUP1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ConsumerGroupID(v.Input)
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
		if actual.IotHubName != v.Expected.IotHubName {
			t.Fatalf("Expected %q but got %q for IotHubName", v.Expected.IotHubName, actual.IotHubName)
		}
		if actual.EventHubEndpointName != v.Expected.EventHubEndpointName {
			t.Fatalf("Expected %q but got %q for EventHubEndpointName", v.Expected.EventHubEndpointName, actual.EventHubEndpointName)
		}
		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}
	}
}

func TestConsumerGroupIDInsensitively(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *ConsumerGroupId
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
			// missing IotHubName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Devices/",
			Error: true,
		},

		{
			// missing value for IotHubName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Devices/iotHubs/",
			Error: true,
		},

		{
			// missing EventHubEndpointName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Devices/iotHubs/hub1/",
			Error: true,
		},

		{
			// missing value for EventHubEndpointName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Devices/iotHubs/hub1/eventHubEndpoints/",
			Error: true,
		},

		{
			// missing Name
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Devices/iotHubs/hub1/eventHubEndpoints/events/",
			Error: true,
		},

		{
			// missing value for Name
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Devices/iotHubs/hub1/eventHubEndpoints/events/consumerGroups/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Devices/iotHubs/hub1/eventHubEndpoints/events/consumerGroups/group1",
			Expected: &ConsumerGroupId{
				SubscriptionId:       "12345678-1234-9876-4563-123456789012",
				ResourceGroup:        "resGroup1",
				IotHubName:           "hub1",
				EventHubEndpointName: "events",
				Name:                 "group1",
			},
		},

		{
			// lower-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Devices/iothubs/hub1/eventhubendpoints/events/consumergroups/group1",
			Expected: &ConsumerGroupId{
				SubscriptionId:       "12345678-1234-9876-4563-123456789012",
				ResourceGroup:        "resGroup1",
				IotHubName:           "hub1",
				EventHubEndpointName: "events",
				Name:                 "group1",
			},
		},

		{
			// upper-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Devices/IOTHUBS/hub1/EVENTHUBENDPOINTS/events/CONSUMERGROUPS/group1",
			Expected: &ConsumerGroupId{
				SubscriptionId:       "12345678-1234-9876-4563-123456789012",
				ResourceGroup:        "resGroup1",
				IotHubName:           "hub1",
				EventHubEndpointName: "events",
				Name:                 "group1",
			},
		},

		{
			// mixed-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Devices/IoThUbS/hub1/EvEnThUbEnDpOiNtS/events/CoNsUmErGrOuPs/group1",
			Expected: &ConsumerGroupId{
				SubscriptionId:       "12345678-1234-9876-4563-123456789012",
				ResourceGroup:        "resGroup1",
				IotHubName:           "hub1",
				EventHubEndpointName: "events",
				Name:                 "group1",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ConsumerGroupIDInsensitively(v.Input)
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
		if actual.IotHubName != v.Expected.IotHubName {
			t.Fatalf("Expected %q but got %q for IotHubName", v.Expected.IotHubName, actual.IotHubName)
		}
		if actual.EventHubEndpointName != v.Expected.EventHubEndpointName {
			t.Fatalf("Expected %q but got %q for EventHubEndpointName", v.Expected.EventHubEndpointName, actual.EventHubEndpointName)
		}
		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}
	}
}
