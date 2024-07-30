package managedclusters

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &MeshRevisionProfileId{}

func TestNewMeshRevisionProfileID(t *testing.T) {
	id := NewMeshRevisionProfileID("12345678-1234-9876-4563-123456789012", "locationValue", "meshRevisionProfileValue")

	if id.SubscriptionId != "12345678-1234-9876-4563-123456789012" {
		t.Fatalf("Expected %q but got %q for Segment 'SubscriptionId'", id.SubscriptionId, "12345678-1234-9876-4563-123456789012")
	}

	if id.LocationName != "locationValue" {
		t.Fatalf("Expected %q but got %q for Segment 'LocationName'", id.LocationName, "locationValue")
	}

	if id.MeshRevisionProfileName != "meshRevisionProfileValue" {
		t.Fatalf("Expected %q but got %q for Segment 'MeshRevisionProfileName'", id.MeshRevisionProfileName, "meshRevisionProfileValue")
	}
}

func TestFormatMeshRevisionProfileID(t *testing.T) {
	actual := NewMeshRevisionProfileID("12345678-1234-9876-4563-123456789012", "locationValue", "meshRevisionProfileValue").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.ContainerService/locations/locationValue/meshRevisionProfiles/meshRevisionProfileValue"
	if actual != expected {
		t.Fatalf("Expected the Formatted ID to be %q but got %q", expected, actual)
	}
}

func TestParseMeshRevisionProfileID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *MeshRevisionProfileId
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
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.ContainerService",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.ContainerService/locations",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.ContainerService/locations/locationValue",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.ContainerService/locations/locationValue/meshRevisionProfiles",
			Error: true,
		},
		{
			// Valid URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.ContainerService/locations/locationValue/meshRevisionProfiles/meshRevisionProfileValue",
			Expected: &MeshRevisionProfileId{
				SubscriptionId:          "12345678-1234-9876-4563-123456789012",
				LocationName:            "locationValue",
				MeshRevisionProfileName: "meshRevisionProfileValue",
			},
		},
		{
			// Invalid (Valid Uri with Extra segment)
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.ContainerService/locations/locationValue/meshRevisionProfiles/meshRevisionProfileValue/extra",
			Error: true,
		},
	}
	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ParseMeshRevisionProfileID(v.Input)
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

		if actual.LocationName != v.Expected.LocationName {
			t.Fatalf("Expected %q but got %q for LocationName", v.Expected.LocationName, actual.LocationName)
		}

		if actual.MeshRevisionProfileName != v.Expected.MeshRevisionProfileName {
			t.Fatalf("Expected %q but got %q for MeshRevisionProfileName", v.Expected.MeshRevisionProfileName, actual.MeshRevisionProfileName)
		}

	}
}

func TestParseMeshRevisionProfileIDInsensitively(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *MeshRevisionProfileId
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
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.ContainerService",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/pRoViDeRs/mIcRoSoFt.cOnTaInErSeRvIcE",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.ContainerService/locations",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/pRoViDeRs/mIcRoSoFt.cOnTaInErSeRvIcE/lOcAtIoNs",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.ContainerService/locations/locationValue",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/pRoViDeRs/mIcRoSoFt.cOnTaInErSeRvIcE/lOcAtIoNs/lOcAtIoNvAlUe",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.ContainerService/locations/locationValue/meshRevisionProfiles",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/pRoViDeRs/mIcRoSoFt.cOnTaInErSeRvIcE/lOcAtIoNs/lOcAtIoNvAlUe/mEsHrEvIsIoNpRoFiLeS",
			Error: true,
		},
		{
			// Valid URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.ContainerService/locations/locationValue/meshRevisionProfiles/meshRevisionProfileValue",
			Expected: &MeshRevisionProfileId{
				SubscriptionId:          "12345678-1234-9876-4563-123456789012",
				LocationName:            "locationValue",
				MeshRevisionProfileName: "meshRevisionProfileValue",
			},
		},
		{
			// Invalid (Valid Uri with Extra segment)
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.ContainerService/locations/locationValue/meshRevisionProfiles/meshRevisionProfileValue/extra",
			Error: true,
		},
		{
			// Valid URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/pRoViDeRs/mIcRoSoFt.cOnTaInErSeRvIcE/lOcAtIoNs/lOcAtIoNvAlUe/mEsHrEvIsIoNpRoFiLeS/mEsHrEvIsIoNpRoFiLeVaLuE",
			Expected: &MeshRevisionProfileId{
				SubscriptionId:          "12345678-1234-9876-4563-123456789012",
				LocationName:            "lOcAtIoNvAlUe",
				MeshRevisionProfileName: "mEsHrEvIsIoNpRoFiLeVaLuE",
			},
		},
		{
			// Invalid (Valid Uri with Extra segment - mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/pRoViDeRs/mIcRoSoFt.cOnTaInErSeRvIcE/lOcAtIoNs/lOcAtIoNvAlUe/mEsHrEvIsIoNpRoFiLeS/mEsHrEvIsIoNpRoFiLeVaLuE/extra",
			Error: true,
		},
	}
	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ParseMeshRevisionProfileIDInsensitively(v.Input)
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

		if actual.LocationName != v.Expected.LocationName {
			t.Fatalf("Expected %q but got %q for LocationName", v.Expected.LocationName, actual.LocationName)
		}

		if actual.MeshRevisionProfileName != v.Expected.MeshRevisionProfileName {
			t.Fatalf("Expected %q but got %q for MeshRevisionProfileName", v.Expected.MeshRevisionProfileName, actual.MeshRevisionProfileName)
		}

	}
}

func TestSegmentsForMeshRevisionProfileId(t *testing.T) {
	segments := MeshRevisionProfileId{}.Segments()
	if len(segments) == 0 {
		t.Fatalf("MeshRevisionProfileId has no segments")
	}

	uniqueNames := make(map[string]struct{}, 0)
	for _, segment := range segments {
		uniqueNames[segment.Name] = struct{}{}
	}
	if len(uniqueNames) != len(segments) {
		t.Fatalf("Expected the Segments to be unique but got %q unique segments and %d total segments", len(uniqueNames), len(segments))
	}
}
