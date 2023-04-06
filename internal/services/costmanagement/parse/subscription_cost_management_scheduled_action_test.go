package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = SubscriptionCostManagementScheduledActionId{}

func TestSubscriptionCostManagementScheduledActionIDFormatter(t *testing.T) {
	actual := NewSubscriptionCostManagementScheduledActionID("12345678-1234-9876-4563-123456789012", "scheduledaction1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.CostManagement/scheduledActions/scheduledaction1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestSubscriptionCostManagementScheduledActionID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *SubscriptionCostManagementScheduledActionId
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
			// missing ScheduledActionName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.CostManagement/",
			Error: true,
		},

		{
			// missing value for ScheduledActionName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.CostManagement/scheduledActions/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.CostManagement/scheduledActions/scheduledaction1",
			Expected: &SubscriptionCostManagementScheduledActionId{
				SubscriptionId:      "12345678-1234-9876-4563-123456789012",
				ScheduledActionName: "scheduledaction1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/PROVIDERS/MICROSOFT.COSTMANAGEMENT/SCHEDULEDACTIONS/SCHEDULEDACTION1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := SubscriptionCostManagementScheduledActionID(v.Input)
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
		if actual.ScheduledActionName != v.Expected.ScheduledActionName {
			t.Fatalf("Expected %q but got %q for ScheduledActionName", v.Expected.ScheduledActionName, actual.ScheduledActionName)
		}
	}
}
