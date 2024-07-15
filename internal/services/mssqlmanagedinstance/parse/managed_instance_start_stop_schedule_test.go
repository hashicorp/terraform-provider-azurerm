// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = ManagedInstanceStartStopScheduleId{}

func TestManagedInstanceStartStopScheduleIDFormatter(t *testing.T) {
	actual := NewManagedInstanceStartStopScheduleID("12345678-1234-9876-4563-123456789012", "resGroup1", "instance1", "default").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Sql/managedInstances/instance1/startStopSchedules/default"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestManagedInstanceStartStopScheduleID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *ManagedInstanceStartStopScheduleId
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
			// missing ManagedInstanceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Sql/",
			Error: true,
		},

		{
			// missing value for ManagedInstanceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Sql/managedInstances/",
			Error: true,
		},

		{
			// missing StartStopScheduleName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Sql/managedInstances/instance1/",
			Error: true,
		},

		{
			// missing value for StartStopScheduleName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Sql/managedInstances/instance1/startStopSchedules/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Sql/managedInstances/instance1/startStopSchedules/default",
			Expected: &ManagedInstanceStartStopScheduleId{
				SubscriptionId:        "12345678-1234-9876-4563-123456789012",
				ResourceGroup:         "resGroup1",
				ManagedInstanceName:   "instance1",
				StartStopScheduleName: "default",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.SQL/MANAGEDINSTANCES/INSTANCE1/STARTSTOPSCHEDULES/DEFAULT",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ManagedInstanceStartStopScheduleID(v.Input)
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
		if actual.ManagedInstanceName != v.Expected.ManagedInstanceName {
			t.Fatalf("Expected %q but got %q for ManagedInstanceName", v.Expected.ManagedInstanceName, actual.ManagedInstanceName)
		}
		if actual.StartStopScheduleName != v.Expected.StartStopScheduleName {
			t.Fatalf("Expected %q but got %q for StartStopScheduleName", v.Expected.StartStopScheduleName, actual.StartStopScheduleName)
		}
	}
}
