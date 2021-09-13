package parse

import (
    "testing"

    "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = DatadogMonitorId {}

func TestDatadogMonitorIDFormatter(t *testing.T) {
    actual := NewDatadogMonitorID("12345678-1234-9876-4563-123456789012", "resourceGroup1", "monitor1").ID()
    expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Datadog/monitors/monitor1"
    if actual != expected {
        t.Fatalf("Expected %q but got %q", expected, actual)
    }
}

func TestDatadogMonitorID(t *testing.T) {
    testData := []struct {
        Input    string
        Error    bool
        Expected *DatadogMonitorId
    }{
        {
            // empty
            Input: "",
            Error: true,
        },

        {
            // missing subscriptions
            Input: "/",
            Error: true,
        },

        {
            // missing value for subscriptions
            Input: "/subscriptions/",
            Error: true,
        },

        {
            // missing resourceGroups
            Input: "/subscriptions/12345678-1234-9876-4563-123456789012/",
            Error: true,
        },

        {
            // missing value for resourceGroups
            Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/",
            Error: true,
        },

        {
            // missing monitors
            Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Datadog/",
            Error: true,
        },

        {
            // missing value for monitors
            Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Datadog/monitors/",
            Error: true,
        },

        {
            // valid
            Input:    "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Datadog/monitors/monitor1",
            Expected: &DatadogMonitorId{
                SubscriptionId:"12345678-1234-9876-4563-123456789012",
                ResourceGroup:"resourceGroup1",
                Name:"monitor1",
            },
        },

        {
            // upper-cased
            Input:    "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESOURCEGROUP1/PROVIDERS/MICROSOFT.DATADOG/MONITORS/MONITOR1",
            Error: true,
        },
    }

    for _, v := range testData {
        t.Logf("[DEBUG] Testing %q", v.Input)

        actual, err := DatadogMonitorID(v.Input)
        if err != nil {
            if v.Error {
                continue
            }

            t.Fatalf("Expected a value but got an error: %s", err)
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

        if actual.Name != v.Expected.Name {
            t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
        }
    }
}