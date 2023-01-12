package validate

import "testing"

func TestNetworkManagerManagementGroupConnectionID(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{

		{
			// empty
			Input: "",
			Valid: false,
		},

		{
			// missing ManagementGroupName
			Input: "/providers/Microsoft.Management/",
			Valid: false,
		},

		{
			// missing value for ManagementGroupName
			Input: "/providers/Microsoft.Management/managementGroups/",
			Valid: false,
		},

		{
			// missing NetworkManagerConnectionName
			Input: "/providers/Microsoft.Management/managementGroups/12345678-1234-9876-4563-123456789012/providers/Microsoft.Network/",
			Valid: false,
		},

		{
			// missing value for NetworkManagerConnectionName
			Input: "/providers/Microsoft.Management/managementGroups/12345678-1234-9876-4563-123456789012/providers/Microsoft.Network/networkManagerConnections/",
			Valid: false,
		},

		{
			// valid
			Input: "/providers/Microsoft.Management/managementGroups/12345678-1234-9876-4563-123456789012/providers/Microsoft.Network/networkManagerConnections/connection1",
			Valid: true,
		},

		{
			// upper-cased
			Input: "/PROVIDERS/MICROSOFT.MANAGEMENT/MANAGEMENTGROUPS/12345678-1234-9876-4563-123456789012/PROVIDERS/MICROSOFT.NETWORK/NETWORKMANAGERCONNECTIONS/CONNECTION1",
			Valid: false,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := NetworkManagerManagementGroupConnectionID(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
