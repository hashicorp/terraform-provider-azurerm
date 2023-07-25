// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = StreamingJobScheduleId{}

func TestStreamingJobScheduleIDFormatter(t *testing.T) {
	actual := NewStreamingJobScheduleID("12345678-1234-9876-4563-123456789012", "resGroup1", "streamingJob1", "default").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.StreamAnalytics/streamingJobs/streamingJob1/schedule/default"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestStreamingJobScheduleID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *StreamingJobScheduleId
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
			// missing StreamingJobName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.StreamAnalytics/",
			Error: true,
		},

		{
			// missing value for StreamingJobName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.StreamAnalytics/streamingJobs/",
			Error: true,
		},

		{
			// missing ScheduleName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.StreamAnalytics/streamingJobs/streamingJob1/",
			Error: true,
		},

		{
			// missing value for ScheduleName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.StreamAnalytics/streamingJobs/streamingJob1/schedule/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.StreamAnalytics/streamingJobs/streamingJob1/schedule/default",
			Expected: &StreamingJobScheduleId{
				SubscriptionId:   "12345678-1234-9876-4563-123456789012",
				ResourceGroup:    "resGroup1",
				StreamingJobName: "streamingJob1",
				ScheduleName:     "default",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.STREAMANALYTICS/STREAMINGJOBS/STREAMINGJOB1/SCHEDULE/DEFAULT",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := StreamingJobScheduleID(v.Input)
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
		if actual.StreamingJobName != v.Expected.StreamingJobName {
			t.Fatalf("Expected %q but got %q for StreamingJobName", v.Expected.StreamingJobName, actual.StreamingJobName)
		}
		if actual.ScheduleName != v.Expected.ScheduleName {
			t.Fatalf("Expected %q but got %q for ScheduleName", v.Expected.ScheduleName, actual.ScheduleName)
		}
	}
}

func TestStreamingJobScheduleIDInsensitively(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *StreamingJobScheduleId
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
			// missing StreamingJobName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.StreamAnalytics/",
			Error: true,
		},

		{
			// missing value for StreamingJobName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.StreamAnalytics/streamingJobs/",
			Error: true,
		},

		{
			// missing ScheduleName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.StreamAnalytics/streamingJobs/streamingJob1/",
			Error: true,
		},

		{
			// missing value for ScheduleName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.StreamAnalytics/streamingJobs/streamingJob1/schedule/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.StreamAnalytics/streamingJobs/streamingJob1/schedule/default",
			Expected: &StreamingJobScheduleId{
				SubscriptionId:   "12345678-1234-9876-4563-123456789012",
				ResourceGroup:    "resGroup1",
				StreamingJobName: "streamingJob1",
				ScheduleName:     "default",
			},
		},

		{
			// lower-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.StreamAnalytics/streamingjobs/streamingJob1/schedule/default",
			Expected: &StreamingJobScheduleId{
				SubscriptionId:   "12345678-1234-9876-4563-123456789012",
				ResourceGroup:    "resGroup1",
				StreamingJobName: "streamingJob1",
				ScheduleName:     "default",
			},
		},

		{
			// upper-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.StreamAnalytics/STREAMINGJOBS/streamingJob1/SCHEDULE/default",
			Expected: &StreamingJobScheduleId{
				SubscriptionId:   "12345678-1234-9876-4563-123456789012",
				ResourceGroup:    "resGroup1",
				StreamingJobName: "streamingJob1",
				ScheduleName:     "default",
			},
		},

		{
			// mixed-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.StreamAnalytics/StReAmInGjObS/streamingJob1/ScHeDuLe/default",
			Expected: &StreamingJobScheduleId{
				SubscriptionId:   "12345678-1234-9876-4563-123456789012",
				ResourceGroup:    "resGroup1",
				StreamingJobName: "streamingJob1",
				ScheduleName:     "default",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := StreamingJobScheduleIDInsensitively(v.Input)
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
		if actual.StreamingJobName != v.Expected.StreamingJobName {
			t.Fatalf("Expected %q but got %q for StreamingJobName", v.Expected.StreamingJobName, actual.StreamingJobName)
		}
		if actual.ScheduleName != v.Expected.ScheduleName {
			t.Fatalf("Expected %q but got %q for ScheduleName", v.Expected.ScheduleName, actual.ScheduleName)
		}
	}
}
