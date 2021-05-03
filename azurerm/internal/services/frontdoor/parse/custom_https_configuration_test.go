package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = CustomHttpsConfigurationId{}

func TestCustomHttpsConfigurationIDFormatter(t *testing.T) {
	actual := NewCustomHttpsConfigurationID("12345678-1234-9876-4563-123456789012", "resGroup1", "frontdoor1", "endpoint1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/frontDoors/frontdoor1/customHttpsConfiguration/endpoint1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestCustomHttpsConfigurationID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *CustomHttpsConfigurationId
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
			// missing FrontDoorName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/",
			Error: true,
		},

		{
			// missing value for FrontDoorName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/frontDoors/",
			Error: true,
		},

		{
			// missing CustomHttpsConfigurationName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/frontDoors/frontdoor1/",
			Error: true,
		},

		{
			// missing value for CustomHttpsConfigurationName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/frontDoors/frontdoor1/customHttpsConfiguration/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/frontDoors/frontdoor1/customHttpsConfiguration/endpoint1",
			Expected: &CustomHttpsConfigurationId{
				SubscriptionId:               "12345678-1234-9876-4563-123456789012",
				ResourceGroup:                "resGroup1",
				FrontDoorName:                "frontdoor1",
				CustomHttpsConfigurationName: "endpoint1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.NETWORK/FRONTDOORS/FRONTDOOR1/CUSTOMHTTPSCONFIGURATION/ENDPOINT1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := CustomHttpsConfigurationID(v.Input)
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
		if actual.FrontDoorName != v.Expected.FrontDoorName {
			t.Fatalf("Expected %q but got %q for FrontDoorName", v.Expected.FrontDoorName, actual.FrontDoorName)
		}
		if actual.CustomHttpsConfigurationName != v.Expected.CustomHttpsConfigurationName {
			t.Fatalf("Expected %q but got %q for CustomHttpsConfigurationName", v.Expected.CustomHttpsConfigurationName, actual.CustomHttpsConfigurationName)
		}
	}
}

func TestCustomHttpsConfigurationIDInsensitively(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *CustomHttpsConfigurationId
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
			// missing FrontDoorName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/",
			Error: true,
		},

		{
			// missing value for FrontDoorName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/frontDoors/",
			Error: true,
		},

		{
			// missing CustomHttpsConfigurationName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/frontDoors/frontdoor1/",
			Error: true,
		},

		{
			// missing value for CustomHttpsConfigurationName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/frontDoors/frontdoor1/customHttpsConfiguration/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/frontDoors/frontdoor1/customHttpsConfiguration/endpoint1",
			Expected: &CustomHttpsConfigurationId{
				SubscriptionId:               "12345678-1234-9876-4563-123456789012",
				ResourceGroup:                "resGroup1",
				FrontDoorName:                "frontdoor1",
				CustomHttpsConfigurationName: "endpoint1",
			},
		},

		{
			// lower-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/frontdoors/frontdoor1/customhttpsconfiguration/endpoint1",
			Expected: &CustomHttpsConfigurationId{
				SubscriptionId:               "12345678-1234-9876-4563-123456789012",
				ResourceGroup:                "resGroup1",
				FrontDoorName:                "frontdoor1",
				CustomHttpsConfigurationName: "endpoint1",
			},
		},

		{
			// upper-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/FRONTDOORS/frontdoor1/CUSTOMHTTPSCONFIGURATION/endpoint1",
			Expected: &CustomHttpsConfigurationId{
				SubscriptionId:               "12345678-1234-9876-4563-123456789012",
				ResourceGroup:                "resGroup1",
				FrontDoorName:                "frontdoor1",
				CustomHttpsConfigurationName: "endpoint1",
			},
		},

		{
			// mixed-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/FrOnTdOoRs/frontdoor1/CuStOmHtTpScOnFiGuRaTiOn/endpoint1",
			Expected: &CustomHttpsConfigurationId{
				SubscriptionId:               "12345678-1234-9876-4563-123456789012",
				ResourceGroup:                "resGroup1",
				FrontDoorName:                "frontdoor1",
				CustomHttpsConfigurationName: "endpoint1",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := CustomHttpsConfigurationIDInsensitively(v.Input)
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
		if actual.FrontDoorName != v.Expected.FrontDoorName {
			t.Fatalf("Expected %q but got %q for FrontDoorName", v.Expected.FrontDoorName, actual.FrontDoorName)
		}
		if actual.CustomHttpsConfigurationName != v.Expected.CustomHttpsConfigurationName {
			t.Fatalf("Expected %q but got %q for CustomHttpsConfigurationName", v.Expected.CustomHttpsConfigurationName, actual.CustomHttpsConfigurationName)
		}
	}
}
