package sourcecontrolconfiguration

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = SourceControlConfigurationId{}

func TestNewSourceControlConfigurationID(t *testing.T) {
	id := NewSourceControlConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterRpValue", "clusterResourceValue", "clusterValue", "sourceControlConfigurationValue")

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

	if id.SourceControlConfigurationName != "sourceControlConfigurationValue" {
		t.Fatalf("Expected %q but got %q for Segment 'SourceControlConfigurationName'", id.SourceControlConfigurationName, "sourceControlConfigurationValue")
	}
}

func TestFormatSourceControlConfigurationID(t *testing.T) {
	actual := NewSourceControlConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterRpValue", "clusterResourceValue", "clusterValue", "sourceControlConfigurationValue").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/clusterRpValue/clusterResourceValue/clusterValue/providers/Microsoft.KubernetesConfiguration/sourceControlConfigurations/sourceControlConfigurationValue"
	if actual != expected {
		t.Fatalf("Expected the Formatted ID to be %q but got %q", expected, actual)
	}
}

func TestParseSourceControlConfigurationID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *SourceControlConfigurationId
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
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/clusterRpValue/clusterResourceValue/clusterValue/providers/Microsoft.KubernetesConfiguration/sourceControlConfigurations",
			Error: true,
		},
		{
			// Valid URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/clusterRpValue/clusterResourceValue/clusterValue/providers/Microsoft.KubernetesConfiguration/sourceControlConfigurations/sourceControlConfigurationValue",
			Expected: &SourceControlConfigurationId{
				SubscriptionId:                 "12345678-1234-9876-4563-123456789012",
				ResourceGroupName:              "example-resource-group",
				ClusterRp:                      "clusterRpValue",
				ClusterResourceName:            "clusterResourceValue",
				ClusterName:                    "clusterValue",
				SourceControlConfigurationName: "sourceControlConfigurationValue",
			},
		},
		{
			// Invalid (Valid Uri with Extra segment)
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/clusterRpValue/clusterResourceValue/clusterValue/providers/Microsoft.KubernetesConfiguration/sourceControlConfigurations/sourceControlConfigurationValue/extra",
			Error: true,
		},
	}
	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ParseSourceControlConfigurationID(v.Input)
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

		if actual.SourceControlConfigurationName != v.Expected.SourceControlConfigurationName {
			t.Fatalf("Expected %q but got %q for SourceControlConfigurationName", v.Expected.SourceControlConfigurationName, actual.SourceControlConfigurationName)
		}

	}
}

func TestParseSourceControlConfigurationIDInsensitively(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *SourceControlConfigurationId
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
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/clusterRpValue/clusterResourceValue/clusterValue/providers/Microsoft.KubernetesConfiguration/sourceControlConfigurations",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/eXaMpLe-rEsOuRcE-GrOuP/pRoViDeRs/cLuStErRpVaLuE/cLuStErReSoUrCeVaLuE/cLuStErVaLuE/pRoViDeRs/mIcRoSoFt.kUbErNeTeScOnFiGuRaTiOn/sOuRcEcOnTrOlCoNfIgUrAtIoNs",
			Error: true,
		},
		{
			// Valid URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/clusterRpValue/clusterResourceValue/clusterValue/providers/Microsoft.KubernetesConfiguration/sourceControlConfigurations/sourceControlConfigurationValue",
			Expected: &SourceControlConfigurationId{
				SubscriptionId:                 "12345678-1234-9876-4563-123456789012",
				ResourceGroupName:              "example-resource-group",
				ClusterRp:                      "clusterRpValue",
				ClusterResourceName:            "clusterResourceValue",
				ClusterName:                    "clusterValue",
				SourceControlConfigurationName: "sourceControlConfigurationValue",
			},
		},
		{
			// Invalid (Valid Uri with Extra segment)
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/clusterRpValue/clusterResourceValue/clusterValue/providers/Microsoft.KubernetesConfiguration/sourceControlConfigurations/sourceControlConfigurationValue/extra",
			Error: true,
		},
		{
			// Valid URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/eXaMpLe-rEsOuRcE-GrOuP/pRoViDeRs/cLuStErRpVaLuE/cLuStErReSoUrCeVaLuE/cLuStErVaLuE/pRoViDeRs/mIcRoSoFt.kUbErNeTeScOnFiGuRaTiOn/sOuRcEcOnTrOlCoNfIgUrAtIoNs/sOuRcEcOnTrOlCoNfIgUrAtIoNvAlUe",
			Expected: &SourceControlConfigurationId{
				SubscriptionId:                 "12345678-1234-9876-4563-123456789012",
				ResourceGroupName:              "eXaMpLe-rEsOuRcE-GrOuP",
				ClusterRp:                      "cLuStErRpVaLuE",
				ClusterResourceName:            "cLuStErReSoUrCeVaLuE",
				ClusterName:                    "cLuStErVaLuE",
				SourceControlConfigurationName: "sOuRcEcOnTrOlCoNfIgUrAtIoNvAlUe",
			},
		},
		{
			// Invalid (Valid Uri with Extra segment - mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/eXaMpLe-rEsOuRcE-GrOuP/pRoViDeRs/cLuStErRpVaLuE/cLuStErReSoUrCeVaLuE/cLuStErVaLuE/pRoViDeRs/mIcRoSoFt.kUbErNeTeScOnFiGuRaTiOn/sOuRcEcOnTrOlCoNfIgUrAtIoNs/sOuRcEcOnTrOlCoNfIgUrAtIoNvAlUe/extra",
			Error: true,
		},
	}
	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ParseSourceControlConfigurationIDInsensitively(v.Input)
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

		if actual.SourceControlConfigurationName != v.Expected.SourceControlConfigurationName {
			t.Fatalf("Expected %q but got %q for SourceControlConfigurationName", v.Expected.SourceControlConfigurationName, actual.SourceControlConfigurationName)
		}

	}
}

func TestSegmentsForSourceControlConfigurationId(t *testing.T) {
	segments := SourceControlConfigurationId{}.Segments()
	if len(segments) == 0 {
		t.Fatalf("SourceControlConfigurationId has no segments")
	}

	uniqueNames := make(map[string]struct{}, 0)
	for _, segment := range segments {
		uniqueNames[segment.Name] = struct{}{}
	}
	if len(uniqueNames) != len(segments) {
		t.Fatalf("Expected the Segments to be unique but got %q unique segments and %d total segments", len(uniqueNames), len(segments))
	}
}
