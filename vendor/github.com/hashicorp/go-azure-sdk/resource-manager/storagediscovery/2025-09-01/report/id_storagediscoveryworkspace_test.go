package report

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &StorageDiscoveryWorkspaceId{}

func TestNewStorageDiscoveryWorkspaceID(t *testing.T) {
	id := NewStorageDiscoveryWorkspaceID("12345678-1234-9876-4563-123456789012", "storageDiscoveryWorkspaceName")

	if id.SubscriptionId != "12345678-1234-9876-4563-123456789012" {
		t.Fatalf("Expected %q but got %q for Segment 'SubscriptionId'", id.SubscriptionId, "12345678-1234-9876-4563-123456789012")
	}

	if id.StorageDiscoveryWorkspaceName != "storageDiscoveryWorkspaceName" {
		t.Fatalf("Expected %q but got %q for Segment 'StorageDiscoveryWorkspaceName'", id.StorageDiscoveryWorkspaceName, "storageDiscoveryWorkspaceName")
	}
}

func TestFormatStorageDiscoveryWorkspaceID(t *testing.T) {
	actual := NewStorageDiscoveryWorkspaceID("12345678-1234-9876-4563-123456789012", "storageDiscoveryWorkspaceName").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.StorageDiscovery/storageDiscoveryWorkspaces/storageDiscoveryWorkspaceName"
	if actual != expected {
		t.Fatalf("Expected the Formatted ID to be %q but got %q", expected, actual)
	}
}

func TestParseStorageDiscoveryWorkspaceID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *StorageDiscoveryWorkspaceId
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
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.StorageDiscovery",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.StorageDiscovery/storageDiscoveryWorkspaces",
			Error: true,
		},
		{
			// Valid URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.StorageDiscovery/storageDiscoveryWorkspaces/storageDiscoveryWorkspaceName",
			Expected: &StorageDiscoveryWorkspaceId{
				SubscriptionId:                "12345678-1234-9876-4563-123456789012",
				StorageDiscoveryWorkspaceName: "storageDiscoveryWorkspaceName",
			},
		},
		{
			// Invalid (Valid Uri with Extra segment)
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.StorageDiscovery/storageDiscoveryWorkspaces/storageDiscoveryWorkspaceName/extra",
			Error: true,
		},
	}
	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ParseStorageDiscoveryWorkspaceID(v.Input)
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

		if actual.StorageDiscoveryWorkspaceName != v.Expected.StorageDiscoveryWorkspaceName {
			t.Fatalf("Expected %q but got %q for StorageDiscoveryWorkspaceName", v.Expected.StorageDiscoveryWorkspaceName, actual.StorageDiscoveryWorkspaceName)
		}

	}
}

func TestParseStorageDiscoveryWorkspaceIDInsensitively(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *StorageDiscoveryWorkspaceId
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
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.StorageDiscovery",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/pRoViDeRs/mIcRoSoFt.sToRaGeDiScOvErY",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.StorageDiscovery/storageDiscoveryWorkspaces",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/pRoViDeRs/mIcRoSoFt.sToRaGeDiScOvErY/sToRaGeDiScOvErYwOrKsPaCeS",
			Error: true,
		},
		{
			// Valid URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.StorageDiscovery/storageDiscoveryWorkspaces/storageDiscoveryWorkspaceName",
			Expected: &StorageDiscoveryWorkspaceId{
				SubscriptionId:                "12345678-1234-9876-4563-123456789012",
				StorageDiscoveryWorkspaceName: "storageDiscoveryWorkspaceName",
			},
		},
		{
			// Invalid (Valid Uri with Extra segment)
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.StorageDiscovery/storageDiscoveryWorkspaces/storageDiscoveryWorkspaceName/extra",
			Error: true,
		},
		{
			// Valid URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/pRoViDeRs/mIcRoSoFt.sToRaGeDiScOvErY/sToRaGeDiScOvErYwOrKsPaCeS/sToRaGeDiScOvErYwOrKsPaCeNaMe",
			Expected: &StorageDiscoveryWorkspaceId{
				SubscriptionId:                "12345678-1234-9876-4563-123456789012",
				StorageDiscoveryWorkspaceName: "sToRaGeDiScOvErYwOrKsPaCeNaMe",
			},
		},
		{
			// Invalid (Valid Uri with Extra segment - mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/pRoViDeRs/mIcRoSoFt.sToRaGeDiScOvErY/sToRaGeDiScOvErYwOrKsPaCeS/sToRaGeDiScOvErYwOrKsPaCeNaMe/extra",
			Error: true,
		},
	}
	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ParseStorageDiscoveryWorkspaceIDInsensitively(v.Input)
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

		if actual.StorageDiscoveryWorkspaceName != v.Expected.StorageDiscoveryWorkspaceName {
			t.Fatalf("Expected %q but got %q for StorageDiscoveryWorkspaceName", v.Expected.StorageDiscoveryWorkspaceName, actual.StorageDiscoveryWorkspaceName)
		}

	}
}

func TestSegmentsForStorageDiscoveryWorkspaceId(t *testing.T) {
	segments := StorageDiscoveryWorkspaceId{}.Segments()
	if len(segments) == 0 {
		t.Fatalf("StorageDiscoveryWorkspaceId has no segments")
	}

	uniqueNames := make(map[string]struct{}, 0)
	for _, segment := range segments {
		uniqueNames[segment.Name] = struct{}{}
	}
	if len(uniqueNames) != len(segments) {
		t.Fatalf("Expected the Segments to be unique but got %q unique segments and %d total segments", len(uniqueNames), len(segments))
	}
}
