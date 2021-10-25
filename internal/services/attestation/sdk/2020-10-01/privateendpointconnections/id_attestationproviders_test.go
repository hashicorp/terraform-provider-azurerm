package privateendpointconnections

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = AttestationProvidersId{}

func TestAttestationProvidersIDFormatter(t *testing.T) {
	actual := NewAttestationProvidersID("{subscriptionId}", "{resourceGroupName}", "{providerName}").ID()
	expected := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Attestation/attestationProviders/{providerName}"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestParseAttestationProvidersID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *AttestationProvidersId
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
			Input: "/subscriptions/{subscriptionId}/",
			Error: true,
		},

		{
			// missing value for ResourceGroup
			Input: "/subscriptions/{subscriptionId}/resourceGroups/",
			Error: true,
		},

		{
			// missing AttestationProviderName
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Attestation/",
			Error: true,
		},

		{
			// missing value for AttestationProviderName
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Attestation/attestationProviders/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Attestation/attestationProviders/{providerName}",
			Expected: &AttestationProvidersId{
				SubscriptionId:          "{subscriptionId}",
				ResourceGroup:           "{resourceGroupName}",
				AttestationProviderName: "{providerName}",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/{SUBSCRIPTIONID}/RESOURCEGROUPS/{RESOURCEGROUPNAME}/PROVIDERS/MICROSOFT.ATTESTATION/ATTESTATIONPROVIDERS/{PROVIDERNAME}",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ParseAttestationProvidersID(v.Input)
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
		if actual.AttestationProviderName != v.Expected.AttestationProviderName {
			t.Fatalf("Expected %q but got %q for AttestationProviderName", v.Expected.AttestationProviderName, actual.AttestationProviderName)
		}
	}
}

func TestParseAttestationProvidersIDInsensitively(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *AttestationProvidersId
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
			Input: "/subscriptions/{subscriptionId}/",
			Error: true,
		},

		{
			// missing value for ResourceGroup
			Input: "/subscriptions/{subscriptionId}/resourceGroups/",
			Error: true,
		},

		{
			// missing AttestationProviderName
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Attestation/",
			Error: true,
		},

		{
			// missing value for AttestationProviderName
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Attestation/attestationProviders/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Attestation/attestationProviders/{providerName}",
			Expected: &AttestationProvidersId{
				SubscriptionId:          "{subscriptionId}",
				ResourceGroup:           "{resourceGroupName}",
				AttestationProviderName: "{providerName}",
			},
		},

		{
			// lower-cased segment names
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Attestation/attestationproviders/{providerName}",
			Expected: &AttestationProvidersId{
				SubscriptionId:          "{subscriptionId}",
				ResourceGroup:           "{resourceGroupName}",
				AttestationProviderName: "{providerName}",
			},
		},

		{
			// upper-cased segment names
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Attestation/ATTESTATIONPROVIDERS/{providerName}",
			Expected: &AttestationProvidersId{
				SubscriptionId:          "{subscriptionId}",
				ResourceGroup:           "{resourceGroupName}",
				AttestationProviderName: "{providerName}",
			},
		},

		{
			// mixed-cased segment names
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Attestation/AtTeStAtIoNpRoViDeRs/{providerName}",
			Expected: &AttestationProvidersId{
				SubscriptionId:          "{subscriptionId}",
				ResourceGroup:           "{resourceGroupName}",
				AttestationProviderName: "{providerName}",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ParseAttestationProvidersIDInsensitively(v.Input)
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
		if actual.AttestationProviderName != v.Expected.AttestationProviderName {
			t.Fatalf("Expected %q but got %q for AttestationProviderName", v.Expected.AttestationProviderName, actual.AttestationProviderName)
		}
	}
}
