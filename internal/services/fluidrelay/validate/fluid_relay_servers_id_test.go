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
			Input: "/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/",
			Valid: false,
		},

		{
			// missing value for ResourceGroup
			Input: "/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/",
			Valid: false,
		},

		{
			// missing FluidRelayServerName
			Input: "/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/myrg/providers/Microsoft.FluidRelay/",
			Valid: false,
		},

		{
			// missing value for FluidRelayServerName
			Input: "/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/myrg/providers/Microsoft.FluidRelay/fluidRelayServers/",
			Valid: false,
		},

		{
			// valid
			Input: "/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/myrg/providers/Microsoft.FluidRelay/fluidRelayServers/myFluid",
			Valid: true,
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/67A9759D-D099-4AA8-8675-E6CFD669C3F4/RESOURCEGROUPS/MYRG/PROVIDERS/MICROSOFT.FLUIDRELAY/FLUIDRELAYSERVERS/MYFLUID",
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
