package validate

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import "testing"

func TestFluidRelayServersID(t *testing.T) {
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
			// missing SubscriptionId
			Input: "/",
			Valid: false,
		},

		{
			// missing value for SubscriptionId
			Input: "/subscriptions/",
			Valid: false,
		},

		{
			// missing ResourceGroup
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/",
			Valid: false,
		},

		{
			// missing value for ResourceGroup
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/",
			Valid: false,
		},

		{
			// missing FluidRelayServerName
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.FluidRelay/",
			Valid: false,
		},

		{
			// missing value for FluidRelayServerName
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.FluidRelay/fluidRelayServers/",
			Valid: false,
		},

		{
			// valid
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.FluidRelay/fluidRelayServers/server1",
			Valid: true,
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/00000000-0000-0000-0000-000000000000/RESOURCEGROUPS/RG1/PROVIDERS/MICROSOFT.FLUIDRELAY/FLUIDRELAYSERVERS/SERVER1",
			Valid: false,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := FluidRelayServersID(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
