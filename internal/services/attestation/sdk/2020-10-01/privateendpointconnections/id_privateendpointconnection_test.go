package privateendpointconnections

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = PrivateEndpointConnectionId{}

func TestPrivateEndpointConnectionIDFormatter(t *testing.T) {
	actual := NewPrivateEndpointConnectionID("{subscriptionId}", "{resourceGroupName}", "{providerName}", "{privateEndpointConnectionName}").ID()
	expected := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Attestation/attestationProviders/{providerName}/privateEndpointConnections/{privateEndpointConnectionName}"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestParsePrivateEndpointConnectionID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *PrivateEndpointConnectionId
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
			// missing Name
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Attestation/attestationProviders/{providerName}/",
			Error: true,
		},

		{
			// missing value for Name
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Attestation/attestationProviders/{providerName}/privateEndpointConnections/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Attestation/attestationProviders/{providerName}/privateEndpointConnections/{privateEndpointConnectionName}",
			Expected: &PrivateEndpointConnectionId{
				SubscriptionId:          "{subscriptionId}",
				ResourceGroup:           "{resourceGroupName}",
				AttestationProviderName: "{providerName}",
				Name:                    "{privateEndpointConnectionName}",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/{SUBSCRIPTIONID}/RESOURCEGROUPS/{RESOURCEGROUPNAME}/PROVIDERS/MICROSOFT.ATTESTATION/ATTESTATIONPROVIDERS/{PROVIDERNAME}/PRIVATEENDPOINTCONNECTIONS/{PRIVATEENDPOINTCONNECTIONNAME}",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ParsePrivateEndpointConnectionID(v.Input)
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
		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}
	}
}

func TestParsePrivateEndpointConnectionIDInsensitively(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *PrivateEndpointConnectionId
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
			// missing Name
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Attestation/attestationProviders/{providerName}/",
			Error: true,
		},

		{
			// missing value for Name
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Attestation/attestationProviders/{providerName}/privateEndpointConnections/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Attestation/attestationProviders/{providerName}/privateEndpointConnections/{privateEndpointConnectionName}",
			Expected: &PrivateEndpointConnectionId{
				SubscriptionId:          "{subscriptionId}",
				ResourceGroup:           "{resourceGroupName}",
				AttestationProviderName: "{providerName}",
				Name:                    "{privateEndpointConnectionName}",
			},
		},

		{
			// lower-cased segment names
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Attestation/attestationproviders/{providerName}/privateendpointconnections/{privateEndpointConnectionName}",
			Expected: &PrivateEndpointConnectionId{
				SubscriptionId:          "{subscriptionId}",
				ResourceGroup:           "{resourceGroupName}",
				AttestationProviderName: "{providerName}",
				Name:                    "{privateEndpointConnectionName}",
			},
		},

		{
			// upper-cased segment names
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Attestation/ATTESTATIONPROVIDERS/{providerName}/PRIVATEENDPOINTCONNECTIONS/{privateEndpointConnectionName}",
			Expected: &PrivateEndpointConnectionId{
				SubscriptionId:          "{subscriptionId}",
				ResourceGroup:           "{resourceGroupName}",
				AttestationProviderName: "{providerName}",
				Name:                    "{privateEndpointConnectionName}",
			},
		},

		{
			// mixed-cased segment names
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Attestation/AtTeStAtIoNpRoViDeRs/{providerName}/PrIvAtEeNdPoInTcOnNeCtIoNs/{privateEndpointConnectionName}",
			Expected: &PrivateEndpointConnectionId{
				SubscriptionId:          "{subscriptionId}",
				ResourceGroup:           "{resourceGroupName}",
				AttestationProviderName: "{providerName}",
				Name:                    "{privateEndpointConnectionName}",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ParsePrivateEndpointConnectionIDInsensitively(v.Input)
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
		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}
	}
}
