// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = NewRelicMonitoredSubscriptionId{}

func TestNewRelicMonitoredSubscriptionIDFormatter(t *testing.T) {
	actual := NewNewRelicMonitoredSubscriptionID("12345678-1234-9876-4563-123456789012", "group1", "monitor1", "default").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/NewRelic.Observability/monitors/monitor1/monitoredSubscriptions/default"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestNewRelicMonitoredSubscriptionID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *NewRelicMonitoredSubscriptionId
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
			// missing MonitorName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/NewRelic.Observability/",
			Error: true,
		},

		{
			// missing value for MonitorName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/NewRelic.Observability/monitors/",
			Error: true,
		},

		{
			// missing MonitoredSubscriptionName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/NewRelic.Observability/monitors/monitor1/",
			Error: true,
		},

		{
			// missing value for MonitoredSubscriptionName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/NewRelic.Observability/monitors/monitor1/monitoredSubscriptions/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/NewRelic.Observability/monitors/monitor1/monitoredSubscriptions/default",
			Expected: &NewRelicMonitoredSubscriptionId{
				SubscriptionId:            "12345678-1234-9876-4563-123456789012",
				ResourceGroup:             "group1",
				MonitorName:               "monitor1",
				MonitoredSubscriptionName: "default",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/GROUP1/PROVIDERS/NEWRELIC.OBSERVABILITY/MONITORS/MONITOR1/MONITOREDSUBSCRIPTIONS/DEFAULT",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := NewRelicMonitoredSubscriptionID(v.Input)
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
		if actual.MonitorName != v.Expected.MonitorName {
			t.Fatalf("Expected %q but got %q for MonitorName", v.Expected.MonitorName, actual.MonitorName)
		}
		if actual.MonitoredSubscriptionName != v.Expected.MonitoredSubscriptionName {
			t.Fatalf("Expected %q but got %q for MonitoredSubscriptionName", v.Expected.MonitoredSubscriptionName, actual.MonitoredSubscriptionName)
		}
	}
}
