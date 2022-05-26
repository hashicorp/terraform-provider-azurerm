package fluxconfigurationoperationstatus

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = FluxConfigurationOperationId{}

func TestNewFluxConfigurationOperationID(t *testing.T) {
	id := NewFluxConfigurationOperationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterRpValue", "clusterResourceValue", "clusterValue", "fluxConfigurationValue", "operationIdValue")

	if id.SubscriptionId != "12345678-1234-9876-4563-123456789012" {
		t.Fatalf("Expected %q but got %q for Segment 'SubscriptionId'", id.SubscriptionId, "12345678-1234-9876-4563-123456789012")
	}

	if id.ResourceGroupName != "example-resource-group" {
		t.Fatalf("Expected %q but got %q for Segment 'ResourceGroupName'", id.ResourceGroupName, "example-resource-group")
	}

	if id.ClusterRp != "clusterRpValue" {
		t.Fatalf("Expected %q but got %q for Segment 'ClusterRp'", id.ClusterRp, "clusterRpValue")
	}

	if id.ClusterResourceName != "clusterResourceValue" {
		t.Fatalf("Expected %q but got %q for Segment 'ClusterResourceName'", id.ClusterResourceName, "clusterResourceValue")
	}

	if id.ClusterName != "clusterValue" {
		t.Fatalf("Expected %q but got %q for Segment 'ClusterName'", id.ClusterName, "clusterValue")
	}

	if id.FluxConfigurationName != "fluxConfigurationValue" {
		t.Fatalf("Expected %q but got %q for Segment 'FluxConfigurationName'", id.FluxConfigurationName, "fluxConfigurationValue")
	}

	if id.OperationId != "operationIdValue" {
		t.Fatalf("Expected %q but got %q for Segment 'OperationId'", id.OperationId, "operationIdValue")
	}
}

func TestFormatFluxConfigurationOperationID(t *testing.T) {
	actual := NewFluxConfigurationOperationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterRpValue", "clusterResourceValue", "clusterValue", "fluxConfigurationValue", "operationIdValue").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/clusterRpValue/clusterResourceValue/clusterValue/providers/Microsoft.KubernetesConfiguration/fluxConfigurations/fluxConfigurationValue/operations/operationIdValue"
	if actual != expected {
		t.Fatalf("Expected the Formatted ID to be %q but got %q", expected, actual)
	}
}

func TestParseFluxConfigurationOperationID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *FluxConfigurationOperationId
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
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/clusterRpValue",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/clusterRpValue/clusterResourceValue",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/clusterRpValue/clusterResourceValue/clusterValue",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/clusterRpValue/clusterResourceValue/clusterValue/providers",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/clusterRpValue/clusterResourceValue/clusterValue/providers/Microsoft.KubernetesConfiguration",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/clusterRpValue/clusterResourceValue/clusterValue/providers/Microsoft.KubernetesConfiguration/fluxConfigurations",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/clusterRpValue/clusterResourceValue/clusterValue/providers/Microsoft.KubernetesConfiguration/fluxConfigurations/fluxConfigurationValue",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/clusterRpValue/clusterResourceValue/clusterValue/providers/Microsoft.KubernetesConfiguration/fluxConfigurations/fluxConfigurationValue/operations",
			Error: true,
		},
		{
			// Valid URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/clusterRpValue/clusterResourceValue/clusterValue/providers/Microsoft.KubernetesConfiguration/fluxConfigurations/fluxConfigurationValue/operations/operationIdValue",
			Expected: &FluxConfigurationOperationId{
				SubscriptionId:        "12345678-1234-9876-4563-123456789012",
				ResourceGroupName:     "example-resource-group",
				ClusterRp:             "clusterRpValue",
				ClusterResourceName:   "clusterResourceValue",
				ClusterName:           "clusterValue",
				FluxConfigurationName: "fluxConfigurationValue",
				OperationId:           "operationIdValue",
			},
		},
		{
			// Invalid (Valid Uri with Extra segment)
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/clusterRpValue/clusterResourceValue/clusterValue/providers/Microsoft.KubernetesConfiguration/fluxConfigurations/fluxConfigurationValue/operations/operationIdValue/extra",
			Error: true,
		},
	}
	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ParseFluxConfigurationOperationID(v.Input)
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

		if actual.ResourceGroupName != v.Expected.ResourceGroupName {
			t.Fatalf("Expected %q but got %q for ResourceGroupName", v.Expected.ResourceGroupName, actual.ResourceGroupName)
		}

		if actual.ClusterRp != v.Expected.ClusterRp {
			t.Fatalf("Expected %q but got %q for ClusterRp", v.Expected.ClusterRp, actual.ClusterRp)
		}

		if actual.ClusterResourceName != v.Expected.ClusterResourceName {
			t.Fatalf("Expected %q but got %q for ClusterResourceName", v.Expected.ClusterResourceName, actual.ClusterResourceName)
		}

		if actual.ClusterName != v.Expected.ClusterName {
			t.Fatalf("Expected %q but got %q for ClusterName", v.Expected.ClusterName, actual.ClusterName)
		}

		if actual.FluxConfigurationName != v.Expected.FluxConfigurationName {
			t.Fatalf("Expected %q but got %q for FluxConfigurationName", v.Expected.FluxConfigurationName, actual.FluxConfigurationName)
		}

		if actual.OperationId != v.Expected.OperationId {
			t.Fatalf("Expected %q but got %q for OperationId", v.Expected.OperationId, actual.OperationId)
		}

	}
}

func TestParseFluxConfigurationOperationIDInsensitively(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *FluxConfigurationOperationId
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
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/clusterRpValue",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/eXaMpLe-rEsOuRcE-GrOuP/pRoViDeRs/cLuStErRpVaLuE",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/clusterRpValue/clusterResourceValue",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/eXaMpLe-rEsOuRcE-GrOuP/pRoViDeRs/cLuStErRpVaLuE/cLuStErReSoUrCeVaLuE",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/clusterRpValue/clusterResourceValue/clusterValue",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/eXaMpLe-rEsOuRcE-GrOuP/pRoViDeRs/cLuStErRpVaLuE/cLuStErReSoUrCeVaLuE/cLuStErVaLuE",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/clusterRpValue/clusterResourceValue/clusterValue/providers",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/eXaMpLe-rEsOuRcE-GrOuP/pRoViDeRs/cLuStErRpVaLuE/cLuStErReSoUrCeVaLuE/cLuStErVaLuE/pRoViDeRs",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/clusterRpValue/clusterResourceValue/clusterValue/providers/Microsoft.KubernetesConfiguration",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/eXaMpLe-rEsOuRcE-GrOuP/pRoViDeRs/cLuStErRpVaLuE/cLuStErReSoUrCeVaLuE/cLuStErVaLuE/pRoViDeRs/mIcRoSoFt.kUbErNeTeScOnFiGuRaTiOn",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/clusterRpValue/clusterResourceValue/clusterValue/providers/Microsoft.KubernetesConfiguration/fluxConfigurations",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/eXaMpLe-rEsOuRcE-GrOuP/pRoViDeRs/cLuStErRpVaLuE/cLuStErReSoUrCeVaLuE/cLuStErVaLuE/pRoViDeRs/mIcRoSoFt.kUbErNeTeScOnFiGuRaTiOn/fLuXcOnFiGuRaTiOnS",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/clusterRpValue/clusterResourceValue/clusterValue/providers/Microsoft.KubernetesConfiguration/fluxConfigurations/fluxConfigurationValue",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/eXaMpLe-rEsOuRcE-GrOuP/pRoViDeRs/cLuStErRpVaLuE/cLuStErReSoUrCeVaLuE/cLuStErVaLuE/pRoViDeRs/mIcRoSoFt.kUbErNeTeScOnFiGuRaTiOn/fLuXcOnFiGuRaTiOnS/fLuXcOnFiGuRaTiOnVaLuE",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/clusterRpValue/clusterResourceValue/clusterValue/providers/Microsoft.KubernetesConfiguration/fluxConfigurations/fluxConfigurationValue/operations",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/eXaMpLe-rEsOuRcE-GrOuP/pRoViDeRs/cLuStErRpVaLuE/cLuStErReSoUrCeVaLuE/cLuStErVaLuE/pRoViDeRs/mIcRoSoFt.kUbErNeTeScOnFiGuRaTiOn/fLuXcOnFiGuRaTiOnS/fLuXcOnFiGuRaTiOnVaLuE/oPeRaTiOnS",
			Error: true,
		},
		{
			// Valid URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/clusterRpValue/clusterResourceValue/clusterValue/providers/Microsoft.KubernetesConfiguration/fluxConfigurations/fluxConfigurationValue/operations/operationIdValue",
			Expected: &FluxConfigurationOperationId{
				SubscriptionId:        "12345678-1234-9876-4563-123456789012",
				ResourceGroupName:     "example-resource-group",
				ClusterRp:             "clusterRpValue",
				ClusterResourceName:   "clusterResourceValue",
				ClusterName:           "clusterValue",
				FluxConfigurationName: "fluxConfigurationValue",
				OperationId:           "operationIdValue",
			},
		},
		{
			// Invalid (Valid Uri with Extra segment)
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/clusterRpValue/clusterResourceValue/clusterValue/providers/Microsoft.KubernetesConfiguration/fluxConfigurations/fluxConfigurationValue/operations/operationIdValue/extra",
			Error: true,
		},
		{
			// Valid URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/eXaMpLe-rEsOuRcE-GrOuP/pRoViDeRs/cLuStErRpVaLuE/cLuStErReSoUrCeVaLuE/cLuStErVaLuE/pRoViDeRs/mIcRoSoFt.kUbErNeTeScOnFiGuRaTiOn/fLuXcOnFiGuRaTiOnS/fLuXcOnFiGuRaTiOnVaLuE/oPeRaTiOnS/oPeRaTiOnIdVaLuE",
			Expected: &FluxConfigurationOperationId{
				SubscriptionId:        "12345678-1234-9876-4563-123456789012",
				ResourceGroupName:     "eXaMpLe-rEsOuRcE-GrOuP",
				ClusterRp:             "cLuStErRpVaLuE",
				ClusterResourceName:   "cLuStErReSoUrCeVaLuE",
				ClusterName:           "cLuStErVaLuE",
				FluxConfigurationName: "fLuXcOnFiGuRaTiOnVaLuE",
				OperationId:           "oPeRaTiOnIdVaLuE",
			},
		},
		{
			// Invalid (Valid Uri with Extra segment - mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/eXaMpLe-rEsOuRcE-GrOuP/pRoViDeRs/cLuStErRpVaLuE/cLuStErReSoUrCeVaLuE/cLuStErVaLuE/pRoViDeRs/mIcRoSoFt.kUbErNeTeScOnFiGuRaTiOn/fLuXcOnFiGuRaTiOnS/fLuXcOnFiGuRaTiOnVaLuE/oPeRaTiOnS/oPeRaTiOnIdVaLuE/extra",
			Error: true,
		},
	}
	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ParseFluxConfigurationOperationIDInsensitively(v.Input)
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

		if actual.ResourceGroupName != v.Expected.ResourceGroupName {
			t.Fatalf("Expected %q but got %q for ResourceGroupName", v.Expected.ResourceGroupName, actual.ResourceGroupName)
		}

		if actual.ClusterRp != v.Expected.ClusterRp {
			t.Fatalf("Expected %q but got %q for ClusterRp", v.Expected.ClusterRp, actual.ClusterRp)
		}

		if actual.ClusterResourceName != v.Expected.ClusterResourceName {
			t.Fatalf("Expected %q but got %q for ClusterResourceName", v.Expected.ClusterResourceName, actual.ClusterResourceName)
		}

		if actual.ClusterName != v.Expected.ClusterName {
			t.Fatalf("Expected %q but got %q for ClusterName", v.Expected.ClusterName, actual.ClusterName)
		}

		if actual.FluxConfigurationName != v.Expected.FluxConfigurationName {
			t.Fatalf("Expected %q but got %q for FluxConfigurationName", v.Expected.FluxConfigurationName, actual.FluxConfigurationName)
		}

		if actual.OperationId != v.Expected.OperationId {
			t.Fatalf("Expected %q but got %q for OperationId", v.Expected.OperationId, actual.OperationId)
		}

	}
}

func TestSegmentsForFluxConfigurationOperationId(t *testing.T) {
	segments := FluxConfigurationOperationId{}.Segments()
	if len(segments) == 0 {
		t.Fatalf("FluxConfigurationOperationId has no segments")
	}

	uniqueNames := make(map[string]struct{}, 0)
	for _, segment := range segments {
		uniqueNames[segment.Name] = struct{}{}
	}
	if len(uniqueNames) != len(segments) {
		t.Fatalf("Expected the Segments to be unique but got %q unique segments and %d total segments", len(uniqueNames), len(segments))
	}
}
