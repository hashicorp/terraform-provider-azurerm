package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = VpnServerConfigurationPolicyGroupId{}

func TestVpnServerConfigurationPolicyGroupIDFormatter(t *testing.T) {
	actual := NewVpnServerConfigurationPolicyGroupID("12345678-1234-9876-4563-123456789012", "resGroup1", "serverConfiguration1", "configurationPolicyGroup1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/vpnServerConfigurations/serverConfiguration1/configurationPolicyGroups/configurationPolicyGroup1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestVpnServerConfigurationPolicyGroupID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *VpnServerConfigurationPolicyGroupId
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
			// missing VpnServerConfigurationName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/",
			Error: true,
		},

		{
			// missing value for VpnServerConfigurationName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/vpnServerConfigurations/",
			Error: true,
		},

		{
			// missing ConfigurationPolicyGroupName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/vpnServerConfigurations/serverConfiguration1/",
			Error: true,
		},

		{
			// missing value for ConfigurationPolicyGroupName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/vpnServerConfigurations/serverConfiguration1/configurationPolicyGroups/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/vpnServerConfigurations/serverConfiguration1/configurationPolicyGroups/configurationPolicyGroup1",
			Expected: &VpnServerConfigurationPolicyGroupId{
				SubscriptionId:               "12345678-1234-9876-4563-123456789012",
				ResourceGroup:                "resGroup1",
				VpnServerConfigurationName:   "serverConfiguration1",
				ConfigurationPolicyGroupName: "configurationPolicyGroup1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.NETWORK/VPNSERVERCONFIGURATIONS/SERVERCONFIGURATION1/CONFIGURATIONPOLICYGROUPS/CONFIGURATIONPOLICYGROUP1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := VpnServerConfigurationPolicyGroupID(v.Input)
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
		if actual.VpnServerConfigurationName != v.Expected.VpnServerConfigurationName {
			t.Fatalf("Expected %q but got %q for VpnServerConfigurationName", v.Expected.VpnServerConfigurationName, actual.VpnServerConfigurationName)
		}
		if actual.ConfigurationPolicyGroupName != v.Expected.ConfigurationPolicyGroupName {
			t.Fatalf("Expected %q but got %q for ConfigurationPolicyGroupName", v.Expected.ConfigurationPolicyGroupName, actual.ConfigurationPolicyGroupName)
		}
	}
}
