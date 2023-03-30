package parse

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = NetworkManagerManagementGroupConnectionId{}

func TestNetworkManagerManagementGroupConnectionIDFormatter(t *testing.T) {
	actual := NewNetworkManagerManagementGroupConnectionID("12345678-1234-9876-4563-123456789012", "connection1").ID()
	expected := "/providers/Microsoft.Management/managementGroups/12345678-1234-9876-4563-123456789012/providers/Microsoft.Network/networkManagerConnections/connection1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestNetworkManagerManagementGroupConnectionID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *NetworkManagerManagementGroupConnectionId
	}{

		{
			// empty
			Input: "",
			Error: true,
		},

		{
			// missing ManagementGroupName
			Input: "/providers/Microsoft.Management/",
			Error: true,
		},

		{
			// missing value for ManagementGroupName
			Input: "/providers/Microsoft.Management/managementGroups/",
			Error: true,
		},

		{
			// missing NetworkManagerConnectionName
			Input: "/providers/Microsoft.Management/managementGroups/12345678-1234-9876-4563-123456789012/providers/Microsoft.Network/",
			Error: true,
		},

		{
			// missing value for NetworkManagerConnectionName
			Input: "/providers/Microsoft.Management/managementGroups/12345678-1234-9876-4563-123456789012/providers/Microsoft.Network/networkManagerConnections/",
			Error: true,
		},

		{
			// valid
			Input: "/providers/Microsoft.Management/managementGroups/12345678-1234-9876-4563-123456789012/providers/Microsoft.Network/networkManagerConnections/connection1",
			Expected: &NetworkManagerManagementGroupConnectionId{
				ManagementGroupName:          "12345678-1234-9876-4563-123456789012",
				NetworkManagerConnectionName: "connection1",
			},
		},

		{
			// upper-cased
			Input: "/PROVIDERS/MICROSOFT.MANAGEMENT/MANAGEMENTGROUPS/12345678-1234-9876-4563-123456789012/PROVIDERS/MICROSOFT.NETWORK/NETWORKMANAGERCONNECTIONS/CONNECTION1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := NetworkManagerManagementGroupConnectionID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %s", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

		if actual.ManagementGroupName != v.Expected.ManagementGroupName {
			t.Fatalf("Expected %q but got %q for ManagementGroupName", v.Expected.ManagementGroupName, actual.ManagementGroupName)
		}
		if actual.NetworkManagerConnectionName != v.Expected.NetworkManagerConnectionName {
			t.Fatalf("Expected %q but got %q for NetworkManagerConnectionName", v.Expected.NetworkManagerConnectionName, actual.NetworkManagerConnectionName)
		}
	}
}
