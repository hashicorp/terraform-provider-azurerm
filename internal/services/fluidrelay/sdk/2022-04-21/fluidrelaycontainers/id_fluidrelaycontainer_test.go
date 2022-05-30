package fluidrelaycontainers

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = FluidRelayContainerId{}

func TestNewFluidRelayContainerID(t *testing.T) {
	id := NewFluidRelayContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fluidRelayServerValue", "fluidRelayContainerValue")

	if id.SubscriptionId != "12345678-1234-9876-4563-123456789012" {
		t.Fatalf("Expected %q but got %q for Segment 'SubscriptionId'", id.SubscriptionId, "12345678-1234-9876-4563-123456789012")
	}

	if id.ResourceGroup != "example-resource-group" {
		t.Fatalf("Expected %q but got %q for Segment 'ResourceGroup'", id.ResourceGroup, "example-resource-group")
	}

	if id.FluidRelayServerName != "fluidRelayServerValue" {
		t.Fatalf("Expected %q but got %q for Segment 'FluidRelayServerName'", id.FluidRelayServerName, "fluidRelayServerValue")
	}

	if id.FluidRelayContainerName != "fluidRelayContainerValue" {
		t.Fatalf("Expected %q but got %q for Segment 'FluidRelayContainerName'", id.FluidRelayContainerName, "fluidRelayContainerValue")
	}
}

func TestFormatFluidRelayContainerID(t *testing.T) {
	actual := NewFluidRelayContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fluidRelayServerValue", "fluidRelayContainerValue").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.FluidRelay/fluidRelayServers/fluidRelayServerValue/fluidRelayContainers/fluidRelayContainerValue"
	if actual != expected {
		t.Fatalf("Expected the Formatted ID to be %q but got %q", expected, actual)
	}
}

func TestParseFluidRelayContainerID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *FluidRelayContainerId
	}{
		{
			// Incomplete URI
			Input: "",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.FluidRelay",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.FluidRelay/fluidRelayServers",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.FluidRelay/fluidRelayServers/fluidRelayServerValue",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.FluidRelay/fluidRelayServers/fluidRelayServerValue/fluidRelayContainers",
			Error: true,
		},
		{
			// Valid URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.FluidRelay/fluidRelayServers/fluidRelayServerValue/fluidRelayContainers/fluidRelayContainerValue",
			Expected: &FluidRelayContainerId{
				SubscriptionId:          "12345678-1234-9876-4563-123456789012",
				ResourceGroup:           "example-resource-group",
				FluidRelayServerName:    "fluidRelayServerValue",
				FluidRelayContainerName: "fluidRelayContainerValue",
			},
		},
		{
			// Invalid (Valid Uri with Extra segment)
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.FluidRelay/fluidRelayServers/fluidRelayServerValue/fluidRelayContainers/fluidRelayContainerValue/extra",
			Error: true,
		},
	}
	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ParseFluidRelayContainerID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %+v", err)
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

		if actual.FluidRelayContainerName != v.Expected.FluidRelayContainerName {
			t.Fatalf("Expected %q but got %q for FluidRelayContainerName", v.Expected.FluidRelayContainerName, actual.FluidRelayContainerName)
		}

	}
}

func TestParseFluidRelayContainerIDInsensitively(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *FluidRelayContainerId
	}{
		{
			// Incomplete URI
			Input: "",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/eXaMpLe-rEsOuRcE-GrOuP",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/eXaMpLe-rEsOuRcE-GrOuP/pRoViDeRs",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.FluidRelay",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/eXaMpLe-rEsOuRcE-GrOuP/pRoViDeRs/mIcRoSoFt.fLuIdReLaY",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.FluidRelay/fluidRelayServers",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/eXaMpLe-rEsOuRcE-GrOuP/pRoViDeRs/mIcRoSoFt.fLuIdReLaY/fLuIdReLaYsErVeRs",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.FluidRelay/fluidRelayServers/fluidRelayServerValue",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/eXaMpLe-rEsOuRcE-GrOuP/pRoViDeRs/mIcRoSoFt.fLuIdReLaY/fLuIdReLaYsErVeRs/fLuIdReLaYsErVeRvAlUe",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.FluidRelay/fluidRelayServers/fluidRelayServerValue/fluidRelayContainers",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/eXaMpLe-rEsOuRcE-GrOuP/pRoViDeRs/mIcRoSoFt.fLuIdReLaY/fLuIdReLaYsErVeRs/fLuIdReLaYsErVeRvAlUe/fLuIdReLaYcOnTaInErS",
			Error: true,
		},
		{
			// Valid URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.FluidRelay/fluidRelayServers/fluidRelayServerValue/fluidRelayContainers/fluidRelayContainerValue",
			Expected: &FluidRelayContainerId{
				SubscriptionId:          "12345678-1234-9876-4563-123456789012",
				ResourceGroup:           "example-resource-group",
				FluidRelayServerName:    "fluidRelayServerValue",
				FluidRelayContainerName: "fluidRelayContainerValue",
			},
		},
		{
			// Invalid (Valid Uri with Extra segment)
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.FluidRelay/fluidRelayServers/fluidRelayServerValue/fluidRelayContainers/fluidRelayContainerValue/extra",
			Error: true,
		},
		{
			// Valid URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/eXaMpLe-rEsOuRcE-GrOuP/pRoViDeRs/mIcRoSoFt.fLuIdReLaY/fLuIdReLaYsErVeRs/fLuIdReLaYsErVeRvAlUe/fLuIdReLaYcOnTaInErS/fLuIdReLaYcOnTaInErVaLuE",
			Expected: &FluidRelayContainerId{
				SubscriptionId:          "12345678-1234-9876-4563-123456789012",
				ResourceGroup:           "eXaMpLe-rEsOuRcE-GrOuP",
				FluidRelayServerName:    "fLuIdReLaYsErVeRvAlUe",
				FluidRelayContainerName: "fLuIdReLaYcOnTaInErVaLuE",
			},
		},
		{
			// Invalid (Valid Uri with Extra segment - mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/eXaMpLe-rEsOuRcE-GrOuP/pRoViDeRs/mIcRoSoFt.fLuIdReLaY/fLuIdReLaYsErVeRs/fLuIdReLaYsErVeRvAlUe/fLuIdReLaYcOnTaInErS/fLuIdReLaYcOnTaInErVaLuE/extra",
			Error: true,
		},
	}
	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ParseFluidRelayContainerIDInsensitively(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %+v", err)
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

		if actual.FluidRelayContainerName != v.Expected.FluidRelayContainerName {
			t.Fatalf("Expected %q but got %q for FluidRelayContainerName", v.Expected.FluidRelayContainerName, actual.FluidRelayContainerName)
		}

	}
}

func TestSegmentsForFluidRelayContainerId(t *testing.T) {
	segments := FluidRelayContainerId{}.Segments()
	if len(segments) == 0 {
		t.Fatalf("FluidRelayContainerId has no segments")
	}

	uniqueNames := make(map[string]struct{}, 0)
	for _, segment := range segments {
		uniqueNames[segment.Name] = struct{}{}
	}
	if len(uniqueNames) != len(segments) {
		t.Fatalf("Expected the Segments to be unique but got %q unique segments and %d total segments", len(uniqueNames), len(segments))
	}
}
