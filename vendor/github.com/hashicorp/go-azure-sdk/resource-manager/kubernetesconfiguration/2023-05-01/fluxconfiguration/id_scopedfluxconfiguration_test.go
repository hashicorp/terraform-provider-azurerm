package fluxconfiguration

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &ScopedFluxConfigurationId{}

func TestNewScopedFluxConfigurationID(t *testing.T) {
	id := NewScopedFluxConfigurationID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "fluxConfigurationValue")

	if id.Scope != "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group" {
		t.Fatalf("Expected %q but got %q for Segment 'Scope'", id.Scope, "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")
	}

	if id.FluxConfigurationName != "fluxConfigurationValue" {
		t.Fatalf("Expected %q but got %q for Segment 'FluxConfigurationName'", id.FluxConfigurationName, "fluxConfigurationValue")
	}
}

func TestFormatScopedFluxConfigurationID(t *testing.T) {
	actual := NewScopedFluxConfigurationID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "fluxConfigurationValue").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group/providers/Microsoft.KubernetesConfiguration/fluxConfigurations/fluxConfigurationValue"
	if actual != expected {
		t.Fatalf("Expected the Formatted ID to be %q but got %q", expected, actual)
	}
}

func TestParseScopedFluxConfigurationID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *ScopedFluxConfigurationId
	}{
		{
			// Incomplete URI
			Input: "",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group/providers",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group/providers/Microsoft.KubernetesConfiguration",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group/providers/Microsoft.KubernetesConfiguration/fluxConfigurations",
			Error: true,
		},
		{
			// Valid URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group/providers/Microsoft.KubernetesConfiguration/fluxConfigurations/fluxConfigurationValue",
			Expected: &ScopedFluxConfigurationId{
				Scope:                 "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group",
				FluxConfigurationName: "fluxConfigurationValue",
			},
		},
		{
			// Invalid (Valid Uri with Extra segment)
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group/providers/Microsoft.KubernetesConfiguration/fluxConfigurations/fluxConfigurationValue/extra",
			Error: true,
		},
	}
	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ParseScopedFluxConfigurationID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %+v", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

		if actual.Scope != v.Expected.Scope {
			t.Fatalf("Expected %q but got %q for Scope", v.Expected.Scope, actual.Scope)
		}

		if actual.FluxConfigurationName != v.Expected.FluxConfigurationName {
			t.Fatalf("Expected %q but got %q for FluxConfigurationName", v.Expected.FluxConfigurationName, actual.FluxConfigurationName)
		}

	}
}

func TestParseScopedFluxConfigurationIDInsensitively(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *ScopedFluxConfigurationId
	}{
		{
			// Incomplete URI
			Input: "",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/sOmE-ReSoUrCe-gRoUp",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group/providers",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/sOmE-ReSoUrCe-gRoUp/pRoViDeRs",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group/providers/Microsoft.KubernetesConfiguration",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/sOmE-ReSoUrCe-gRoUp/pRoViDeRs/mIcRoSoFt.kUbErNeTeScOnFiGuRaTiOn",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group/providers/Microsoft.KubernetesConfiguration/fluxConfigurations",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/sOmE-ReSoUrCe-gRoUp/pRoViDeRs/mIcRoSoFt.kUbErNeTeScOnFiGuRaTiOn/fLuXcOnFiGuRaTiOnS",
			Error: true,
		},
		{
			// Valid URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group/providers/Microsoft.KubernetesConfiguration/fluxConfigurations/fluxConfigurationValue",
			Expected: &ScopedFluxConfigurationId{
				Scope:                 "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group",
				FluxConfigurationName: "fluxConfigurationValue",
			},
		},
		{
			// Invalid (Valid Uri with Extra segment)
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group/providers/Microsoft.KubernetesConfiguration/fluxConfigurations/fluxConfigurationValue/extra",
			Error: true,
		},
		{
			// Valid URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/sOmE-ReSoUrCe-gRoUp/pRoViDeRs/mIcRoSoFt.kUbErNeTeScOnFiGuRaTiOn/fLuXcOnFiGuRaTiOnS/fLuXcOnFiGuRaTiOnVaLuE",
			Expected: &ScopedFluxConfigurationId{
				Scope:                 "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/sOmE-ReSoUrCe-gRoUp",
				FluxConfigurationName: "fLuXcOnFiGuRaTiOnVaLuE",
			},
		},
		{
			// Invalid (Valid Uri with Extra segment - mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/sOmE-ReSoUrCe-gRoUp/pRoViDeRs/mIcRoSoFt.kUbErNeTeScOnFiGuRaTiOn/fLuXcOnFiGuRaTiOnS/fLuXcOnFiGuRaTiOnVaLuE/extra",
			Error: true,
		},
	}
	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ParseScopedFluxConfigurationIDInsensitively(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %+v", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

		if actual.Scope != v.Expected.Scope {
			t.Fatalf("Expected %q but got %q for Scope", v.Expected.Scope, actual.Scope)
		}

		if actual.FluxConfigurationName != v.Expected.FluxConfigurationName {
			t.Fatalf("Expected %q but got %q for FluxConfigurationName", v.Expected.FluxConfigurationName, actual.FluxConfigurationName)
		}

	}
}

func TestSegmentsForScopedFluxConfigurationId(t *testing.T) {
	segments := ScopedFluxConfigurationId{}.Segments()
	if len(segments) == 0 {
		t.Fatalf("ScopedFluxConfigurationId has no segments")
	}

	uniqueNames := make(map[string]struct{}, 0)
	for _, segment := range segments {
		uniqueNames[segment.Name] = struct{}{}
	}
	if len(uniqueNames) != len(segments) {
		t.Fatalf("Expected the Segments to be unique but got %q unique segments and %d total segments", len(uniqueNames), len(segments))
	}
}
