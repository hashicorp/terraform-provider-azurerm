package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = BackendAddressPoolAddressId{}

func TestBackendAddressPoolAddressIDFormatter(t *testing.T) {
	actual := NewBackendAddressPoolAddressID("12345678-1234-9876-4563-123456789012", "resGroup1", "loadBalancer1", "backendAddressPool1", "address1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/loadBalancers/loadBalancer1/backendAddressPools/backendAddressPool1/addresses/address1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestBackendAddressPoolAddressID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *BackendAddressPoolAddressId
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
			// missing AddressName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/loadBalancers/loadBalancer1/backendAddressPools/backendAddressPool1/",
			Error: true,
		},

		{
			// missing value for AddressName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/loadBalancers/loadBalancer1/backendAddressPools/backendAddressPool1/addresses/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/loadBalancers/loadBalancer1/backendAddressPools/backendAddressPool1/addresses/address1",
			Expected: &BackendAddressPoolAddressId{
				SubscriptionId:         "12345678-1234-9876-4563-123456789012",
				ResourceGroup:          "resGroup1",
				LoadBalancerName:       "loadBalancer1",
				BackendAddressPoolName: "backendAddressPool1",
				AddressName:            "address1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.NETWORK/LOADBALANCERS/LOADBALANCER1/BACKENDADDRESSPOOLS/BACKENDADDRESSPOOL1/ADDRESSES/ADDRESS1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := BackendAddressPoolAddressID(v.Input)
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
		if actual.AddressName != v.Expected.AddressName {
			t.Fatalf("Expected %q but got %q for AddressName", v.Expected.AddressName, actual.AddressName)
		}
	}
}
