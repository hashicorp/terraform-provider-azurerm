package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = StreamInputId{}

func TestStreamInputIDFormatter(t *testing.T) {
	actual := NewStreamInputID("12345678-1234-9876-4563-123456789012", "resGroup1", "streamingJob1", "streamInput1").ID("")
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.StreamAnalytics/streamingjobs/streamingJob1/inputs/streamInput1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestStreamInputID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *StreamInputId
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
			// missing StreamingjobName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.StreamAnalytics/",
			Error: true,
		},

		{
			// missing value for StreamingjobName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.StreamAnalytics/streamingjobs/",
			Error: true,
		},

		{
			// missing InputName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.StreamAnalytics/streamingjobs/streamingJob1/",
			Error: true,
		},

		{
			// missing value for InputName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.StreamAnalytics/streamingjobs/streamingJob1/inputs/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.StreamAnalytics/streamingjobs/streamingJob1/inputs/streamInput1",
			Expected: &StreamInputId{
				SubscriptionId:   "12345678-1234-9876-4563-123456789012",
				ResourceGroup:    "resGroup1",
				StreamingjobName: "streamingJob1",
				InputName:        "streamInput1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.STREAMANALYTICS/STREAMINGJOBS/STREAMINGJOB1/INPUTS/STREAMINPUT1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := StreamInputID(v.Input)
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
		if actual.StreamingjobName != v.Expected.StreamingjobName {
			t.Fatalf("Expected %q but got %q for StreamingjobName", v.Expected.StreamingjobName, actual.StreamingjobName)
		}
		if actual.InputName != v.Expected.InputName {
			t.Fatalf("Expected %q but got %q for InputName", v.Expected.InputName, actual.InputName)
		}
	}
}
