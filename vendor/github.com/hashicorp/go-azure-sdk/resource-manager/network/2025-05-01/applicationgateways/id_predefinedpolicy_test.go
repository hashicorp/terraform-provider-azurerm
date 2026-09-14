package applicationgateways

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &PredefinedPolicyId{}

func TestNewPredefinedPolicyID(t *testing.T) {
	id := NewPredefinedPolicyID("12345678-1234-9876-4563-123456789012", "predefinedPolicyName")

	if id.SubscriptionId != "12345678-1234-9876-4563-123456789012" {
		t.Fatalf("Expected %q but got %q for Segment 'SubscriptionId'", id.SubscriptionId, "12345678-1234-9876-4563-123456789012")
	}

	if id.PredefinedPolicyName != "predefinedPolicyName" {
		t.Fatalf("Expected %q but got %q for Segment 'PredefinedPolicyName'", id.PredefinedPolicyName, "predefinedPolicyName")
	}
}

func TestFormatPredefinedPolicyID(t *testing.T) {
	actual := NewPredefinedPolicyID("12345678-1234-9876-4563-123456789012", "predefinedPolicyName").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Network/applicationGatewayAvailableSslOptions/default/predefinedPolicies/predefinedPolicyName"
	if actual != expected {
		t.Fatalf("Expected the Formatted ID to be %q but got %q", expected, actual)
	}
}

func TestParsePredefinedPolicyID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *PredefinedPolicyId
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
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Network",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Network/applicationGatewayAvailableSslOptions",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Network/applicationGatewayAvailableSslOptions/default",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Network/applicationGatewayAvailableSslOptions/default/predefinedPolicies",
			Error: true,
		},
		{
			// Valid URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Network/applicationGatewayAvailableSslOptions/default/predefinedPolicies/predefinedPolicyName",
			Expected: &PredefinedPolicyId{
				SubscriptionId:       "12345678-1234-9876-4563-123456789012",
				PredefinedPolicyName: "predefinedPolicyName",
			},
		},
		{
			// Invalid (Valid Uri with Extra segment)
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Network/applicationGatewayAvailableSslOptions/default/predefinedPolicies/predefinedPolicyName/extra",
			Error: true,
		},
	}
	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ParsePredefinedPolicyID(v.Input)
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

		if actual.PredefinedPolicyName != v.Expected.PredefinedPolicyName {
			t.Fatalf("Expected %q but got %q for PredefinedPolicyName", v.Expected.PredefinedPolicyName, actual.PredefinedPolicyName)
		}

	}
}

func TestParsePredefinedPolicyIDInsensitively(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *PredefinedPolicyId
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
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Network",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/pRoViDeRs/mIcRoSoFt.nEtWoRk",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Network/applicationGatewayAvailableSslOptions",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/pRoViDeRs/mIcRoSoFt.nEtWoRk/aPpLiCaTiOnGaTeWaYaVaIlAbLeSsLoPtIoNs",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Network/applicationGatewayAvailableSslOptions/default",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/pRoViDeRs/mIcRoSoFt.nEtWoRk/aPpLiCaTiOnGaTeWaYaVaIlAbLeSsLoPtIoNs/dEfAuLt",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Network/applicationGatewayAvailableSslOptions/default/predefinedPolicies",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/pRoViDeRs/mIcRoSoFt.nEtWoRk/aPpLiCaTiOnGaTeWaYaVaIlAbLeSsLoPtIoNs/dEfAuLt/pReDeFiNeDpOlIcIeS",
			Error: true,
		},
		{
			// Valid URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Network/applicationGatewayAvailableSslOptions/default/predefinedPolicies/predefinedPolicyName",
			Expected: &PredefinedPolicyId{
				SubscriptionId:       "12345678-1234-9876-4563-123456789012",
				PredefinedPolicyName: "predefinedPolicyName",
			},
		},
		{
			// Invalid (Valid Uri with Extra segment)
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Network/applicationGatewayAvailableSslOptions/default/predefinedPolicies/predefinedPolicyName/extra",
			Error: true,
		},
		{
			// Valid URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/pRoViDeRs/mIcRoSoFt.nEtWoRk/aPpLiCaTiOnGaTeWaYaVaIlAbLeSsLoPtIoNs/dEfAuLt/pReDeFiNeDpOlIcIeS/pReDeFiNeDpOlIcYnAmE",
			Expected: &PredefinedPolicyId{
				SubscriptionId:       "12345678-1234-9876-4563-123456789012",
				PredefinedPolicyName: "pReDeFiNeDpOlIcYnAmE",
			},
		},
		{
			// Invalid (Valid Uri with Extra segment - mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/pRoViDeRs/mIcRoSoFt.nEtWoRk/aPpLiCaTiOnGaTeWaYaVaIlAbLeSsLoPtIoNs/dEfAuLt/pReDeFiNeDpOlIcIeS/pReDeFiNeDpOlIcYnAmE/extra",
			Error: true,
		},
	}
	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ParsePredefinedPolicyIDInsensitively(v.Input)
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

		if actual.PredefinedPolicyName != v.Expected.PredefinedPolicyName {
			t.Fatalf("Expected %q but got %q for PredefinedPolicyName", v.Expected.PredefinedPolicyName, actual.PredefinedPolicyName)
		}

	}
}

func TestSegmentsForPredefinedPolicyId(t *testing.T) {
	segments := PredefinedPolicyId{}.Segments()
	if len(segments) == 0 {
		t.Fatalf("PredefinedPolicyId has no segments")
	}

	uniqueNames := make(map[string]struct{}, 0)
	for _, segment := range segments {
		uniqueNames[segment.Name] = struct{}{}
	}
	if len(uniqueNames) != len(segments) {
		t.Fatalf("Expected the Segments to be unique but got %d unique segments and %d total segments", len(uniqueNames), len(segments))
	}
}
