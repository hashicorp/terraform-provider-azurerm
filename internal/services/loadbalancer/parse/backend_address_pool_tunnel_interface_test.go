package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/resourceid"
)

var _ resourceid.Formatter = BackendAddressPoolTunnelInterfaceId{}

func TestBackendAddressPoolTunnelInterfaceIDFormatter(t *testing.T) {
	actual := NewBackendAddressPoolTunnelInterfaceID("12345678-1234-9876-4563-123456789012", "resGroup1", "loadBalancer1", "backendAddressPool1", "itf1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/loadBalancers/loadBalancer1/backendAddressPools/backendAddressPool1/tunnelInterfaces/itf1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestBackendAddressPoolTunnelInterfaceID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *BackendAddressPoolTunnelInterfaceId
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
			// missing LoadBalancerName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/",
			Error: true,
		},

		{
			// missing value for LoadBalancerName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/loadBalancers/",
			Error: true,
		},

		{
			// missing BackendAddressPoolName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/loadBalancers/loadBalancer1/",
			Error: true,
		},

		{
			// missing value for BackendAddressPoolName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/loadBalancers/loadBalancer1/backendAddressPools/",
			Error: true,
		},

		{
			// missing TunnelInterfaceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/loadBalancers/loadBalancer1/backendAddressPools/backendAddressPool1/",
			Error: true,
		},

		{
			// missing value for TunnelInterfaceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/loadBalancers/loadBalancer1/backendAddressPools/backendAddressPool1/tunnelInterfaces/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/loadBalancers/loadBalancer1/backendAddressPools/backendAddressPool1/tunnelInterfaces/itf1",
			Expected: &BackendAddressPoolTunnelInterfaceId{
				SubscriptionId:         "12345678-1234-9876-4563-123456789012",
				ResourceGroup:          "resGroup1",
				LoadBalancerName:       "loadBalancer1",
				BackendAddressPoolName: "backendAddressPool1",
				TunnelInterfaceName:    "itf1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.NETWORK/LOADBALANCERS/LOADBALANCER1/BACKENDADDRESSPOOLS/BACKENDADDRESSPOOL1/TUNNELINTERFACES/ITF1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := BackendAddressPoolTunnelInterfaceID(v.Input)
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
		if actual.LoadBalancerName != v.Expected.LoadBalancerName {
			t.Fatalf("Expected %q but got %q for LoadBalancerName", v.Expected.LoadBalancerName, actual.LoadBalancerName)
		}
		if actual.BackendAddressPoolName != v.Expected.BackendAddressPoolName {
			t.Fatalf("Expected %q but got %q for BackendAddressPoolName", v.Expected.BackendAddressPoolName, actual.BackendAddressPoolName)
		}
		if actual.TunnelInterfaceName != v.Expected.TunnelInterfaceName {
			t.Fatalf("Expected %q but got %q for TunnelInterfaceName", v.Expected.TunnelInterfaceName, actual.TunnelInterfaceName)
		}
	}
}
