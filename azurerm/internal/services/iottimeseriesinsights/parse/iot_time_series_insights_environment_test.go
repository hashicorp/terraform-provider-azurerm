package parse

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = TimeSeriesInsightsEnvironmentId{}

func TestTimeSeriesInsightsEnvironmentIDFormatter(t *testing.T) {
	subscriptionId := "12345678-1234-5678-1234-123456789012"
	actual := NewTimeSeriesInsightsEnvironmentID(subscriptionId, "resourceGroup1", "env1").ID("")
	expected := "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.TimeSeriesInsights/environments/env1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestTimeSeriesInsightsEnvironmentId(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *TimeSeriesInsightsEnvironmentId
	}{
		{
			Name:     "Empty",
			Input:    "",
			Expected: nil,
		},
		{
			Name:     "No Resource Groups Segment",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000",
			Expected: nil,
		},
		{
			Name:     "No Resource Groups Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/",
			Expected: nil,
		},
		{
			Name:     "Resource Group ID",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/",
			Expected: nil,
		},
		{
			Name:     "Time Series Insight Environment Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.TimeSeriesInsights/environments/",
			Expected: nil,
		},
		{
			Name:  "Time Series Insight Environment ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.TimeSeriesInsights/environments/Environment1",
			Expected: &TimeSeriesInsightsEnvironmentId{
				SubscriptionId: "00000000-0000-0000-0000-000000000000",
				Name:           "Environment1",
				ResourceGroup:  "resGroup1",
			},
		},
		{
			Name:     "Wrong Casing",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.TimeSeriesInsights/Environments/Environment1",
			Expected: nil,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := TimeSeriesInsightsEnvironmentID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expected.ResourceGroup, actual.ResourceGroup)
		}

		if actual.SubscriptionId != v.Expected.SubscriptionId {
			t.Fatalf("Expected %q but got %q for Subscription Id", v.Expected.SubscriptionId, actual.SubscriptionId)
		}
	}
}
