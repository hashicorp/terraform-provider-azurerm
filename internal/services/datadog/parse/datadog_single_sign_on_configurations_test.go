package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = DatadogSingleSignOnConfigurationsId{}

func TestDatadogSingleSignOnConfigurationsIDFormatter(t *testing.T) {
	actual := NewDatadogSingleSignOnConfigurationsID("12345678-1234-9876-4563-123456789012", "resourceGroup1", "monitor1", "default").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Datadog/monitors/monitor1/singleSignOnConfigurations/default"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestDatadogSingleSignOnConfigurationsID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *DatadogSingleSignOnConfigurationsId
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
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Datadog/",
			Error: true,
		},

		{
			// missing value for MonitorName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Datadog/monitors/",
			Error: true,
		},

		{
			// missing SingleSignOnConfigurationName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Datadog/monitors/monitor1/",
			Error: true,
		},

		{
			// missing value for SingleSignOnConfigurationName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Datadog/monitors/monitor1/singleSignOnConfigurations/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Datadog/monitors/monitor1/singleSignOnConfigurations/default",
			Expected: &DatadogSingleSignOnConfigurationsId{
				SubscriptionId:                "12345678-1234-9876-4563-123456789012",
				ResourceGroup:                 "resourceGroup1",
				MonitorName:                   "monitor1",
				SingleSignOnConfigurationName: "default",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESOURCEGROUP1/PROVIDERS/MICROSOFT.DATADOG/MONITORS/MONITOR1/SINGLESIGNONCONFIGURATIONS/DEFAULT",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := DatadogSingleSignOnConfigurationsID(v.Input)
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
		if actual.SingleSignOnConfigurationName != v.Expected.SingleSignOnConfigurationName {
			t.Fatalf("Expected %q but got %q for SingleSignOnConfigurationName", v.Expected.SingleSignOnConfigurationName, actual.SingleSignOnConfigurationName)
		}
	}
}
