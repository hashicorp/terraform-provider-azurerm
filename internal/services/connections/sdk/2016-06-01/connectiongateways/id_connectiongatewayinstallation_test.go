package connectiongateways

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = ConnectionGatewayInstallationId{}

func TestNewConnectionGatewayInstallationID(t *testing.T) {
	id := NewConnectionGatewayInstallationID("12345678-1234-9876-4563-123456789012", "locationValue", "gatewayIdValue")

	if id.SubscriptionId != "12345678-1234-9876-4563-123456789012" {
		t.Fatalf("Expected %q but got %q for Segment 'SubscriptionId'", id.SubscriptionId, "12345678-1234-9876-4563-123456789012")
	}

	if id.Location != "locationValue" {
		t.Fatalf("Expected %q but got %q for Segment 'Location'", id.Location, "locationValue")
	}

	if id.GatewayId != "gatewayIdValue" {
		t.Fatalf("Expected %q but got %q for Segment 'GatewayId'", id.GatewayId, "gatewayIdValue")
	}
}

func TestFormatConnectionGatewayInstallationID(t *testing.T) {
	actual := NewConnectionGatewayInstallationID("12345678-1234-9876-4563-123456789012", "locationValue", "gatewayIdValue").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Web/locations/locationValue/connectionGatewayInstallations/gatewayIdValue"
	if actual != expected {
		t.Fatalf("Expected the Formatted ID to be %q but got %q", expected, actual)
	}
}

func TestParseConnectionGatewayInstallationID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *ConnectionGatewayInstallationId
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
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Web",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Web/locations",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Web/locations/locationValue",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Web/locations/locationValue/connectionGatewayInstallations",
			Error: true,
		},
		{
			// Valid URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Web/locations/locationValue/connectionGatewayInstallations/gatewayIdValue",
			Expected: &ConnectionGatewayInstallationId{
				SubscriptionId: "12345678-1234-9876-4563-123456789012",
				Location:       "locationValue",
				GatewayId:      "gatewayIdValue",
			},
		},
		{
			// Invalid (Valid Uri with Extra segment)
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Web/locations/locationValue/connectionGatewayInstallations/gatewayIdValue/extra",
			Error: true,
		},
	}
	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ParseConnectionGatewayInstallationID(v.Input)
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

		if actual.Location != v.Expected.Location {
			t.Fatalf("Expected %q but got %q for Location", v.Expected.Location, actual.Location)
		}

		if actual.GatewayId != v.Expected.GatewayId {
			t.Fatalf("Expected %q but got %q for GatewayId", v.Expected.GatewayId, actual.GatewayId)
		}

	}
}

func TestParseConnectionGatewayInstallationIDInsensitively(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *ConnectionGatewayInstallationId
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
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/pRoViDeRs",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Web",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/pRoViDeRs/mIcRoSoFt.wEb",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Web/locations",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/pRoViDeRs/mIcRoSoFt.wEb/lOcAtIoNs",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Web/locations/locationValue",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/pRoViDeRs/mIcRoSoFt.wEb/lOcAtIoNs/lOcAtIoNvAlUe",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Web/locations/locationValue/connectionGatewayInstallations",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/pRoViDeRs/mIcRoSoFt.wEb/lOcAtIoNs/lOcAtIoNvAlUe/cOnNeCtIoNgAtEwAyInStAlLaTiOnS",
			Error: true,
		},
		{
			// Valid URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Web/locations/locationValue/connectionGatewayInstallations/gatewayIdValue",
			Expected: &ConnectionGatewayInstallationId{
				SubscriptionId: "12345678-1234-9876-4563-123456789012",
				Location:       "locationValue",
				GatewayId:      "gatewayIdValue",
			},
		},
		{
			// Invalid (Valid Uri with Extra segment)
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Web/locations/locationValue/connectionGatewayInstallations/gatewayIdValue/extra",
			Error: true,
		},
		{
			// Valid URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/pRoViDeRs/mIcRoSoFt.wEb/lOcAtIoNs/lOcAtIoNvAlUe/cOnNeCtIoNgAtEwAyInStAlLaTiOnS/gAtEwAyIdVaLuE",
			Expected: &ConnectionGatewayInstallationId{
				SubscriptionId: "12345678-1234-9876-4563-123456789012",
				Location:       "lOcAtIoNvAlUe",
				GatewayId:      "gAtEwAyIdVaLuE",
			},
		},
		{
			// Invalid (Valid Uri with Extra segment - mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/pRoViDeRs/mIcRoSoFt.wEb/lOcAtIoNs/lOcAtIoNvAlUe/cOnNeCtIoNgAtEwAyInStAlLaTiOnS/gAtEwAyIdVaLuE/extra",
			Error: true,
		},
	}
	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ParseConnectionGatewayInstallationIDInsensitively(v.Input)
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

		if actual.Location != v.Expected.Location {
			t.Fatalf("Expected %q but got %q for Location", v.Expected.Location, actual.Location)
		}

		if actual.GatewayId != v.Expected.GatewayId {
			t.Fatalf("Expected %q but got %q for GatewayId", v.Expected.GatewayId, actual.GatewayId)
		}

	}
}

func TestSegmentsForConnectionGatewayInstallationId(t *testing.T) {
	segments := ConnectionGatewayInstallationId{}.Segments()
	if len(segments) == 0 {
		t.Fatalf("ConnectionGatewayInstallationId has no segments")
	}

	uniqueNames := make(map[string]struct{}, 0)
	for _, segment := range segments {
		uniqueNames[segment.Name] = struct{}{}
	}
	if len(uniqueNames) != len(segments) {
		t.Fatalf("Expected the Segments to be unique but got %q unique segments and %d total segments", len(uniqueNames), len(segments))
	}
}
