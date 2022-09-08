package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = FluidRelayId{}

func TestFluidRelayIDFormatter(t *testing.T) {
	actual := NewFluidRelayID("67a9759d-d099-4aa8-8675-e6cfd669c3f4", "myrg", "myFluid").ID()
	expected := "/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/myrg/providers/Microsoft.FluidRelay/fluidRelayServers/myFluid"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestFluidRelayID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *FluidRelayId
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
			Input: "/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/",
			Error: true,
		},

		{
			// missing value for ResourceGroup
			Input: "/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/",
			Error: true,
		},

		{
			// missing FluidRelayServerName
			Input: "/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/myrg/providers/Microsoft.FluidRelay/",
			Error: true,
		},

		{
			// missing value for FluidRelayServerName
			Input: "/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/myrg/providers/Microsoft.FluidRelay/fluidRelayServers/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/myrg/providers/Microsoft.FluidRelay/fluidRelayServers/myFluid",
			Expected: &FluidRelayId{
				SubscriptionId:       "67a9759d-d099-4aa8-8675-e6cfd669c3f4",
				ResourceGroup:        "myrg",
				FluidRelayServerName: "myFluid",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/67A9759D-D099-4AA8-8675-E6CFD669C3F4/RESOURCEGROUPS/MYRG/PROVIDERS/MICROSOFT.FLUIDRELAY/FLUIDRELAYSERVERS/MYFLUID",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := FluidRelayID(v.Input)
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
		if actual.FluidRelayServerName != v.Expected.FluidRelayServerName {
			t.Fatalf("Expected %q but got %q for FluidRelayServerName", v.Expected.FluidRelayServerName, actual.FluidRelayServerName)
		}
	}
}
