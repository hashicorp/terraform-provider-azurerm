package mhsmprivateendpointconnections

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = ManagedHSMPrivateEndpointConnectionId{}

func TestManagedHSMPrivateEndpointConnectionIDFormatter(t *testing.T) {
	actual := NewManagedHSMPrivateEndpointConnectionID("{subscriptionId}", "{resourceGroupName}", "{name}", "{privateEndpointConnectionName}").ID()
	expected := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.KeyVault/managedHSMs/{name}/privateEndpointConnections/{privateEndpointConnectionName}"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestParseManagedHSMPrivateEndpointConnectionID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *ManagedHSMPrivateEndpointConnectionId
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
			// missing ManagedHSMName
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.KeyVault/",
			Error: true,
		},

		{
			// missing value for ManagedHSMName
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.KeyVault/managedHSMs/",
			Error: true,
		},

		{
			// missing PrivateEndpointConnectionName
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.KeyVault/managedHSMs/{name}/",
			Error: true,
		},

		{
			// missing value for PrivateEndpointConnectionName
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.KeyVault/managedHSMs/{name}/privateEndpointConnections/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.KeyVault/managedHSMs/{name}/privateEndpointConnections/{privateEndpointConnectionName}",
			Expected: &ManagedHSMPrivateEndpointConnectionId{
				SubscriptionId:                "{subscriptionId}",
				ResourceGroup:                 "{resourceGroupName}",
				ManagedHSMName:                "{name}",
				PrivateEndpointConnectionName: "{privateEndpointConnectionName}",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/{SUBSCRIPTIONID}/RESOURCEGROUPS/{RESOURCEGROUPNAME}/PROVIDERS/MICROSOFT.KEYVAULT/MANAGEDHSMS/{NAME}/PRIVATEENDPOINTCONNECTIONS/{PRIVATEENDPOINTCONNECTIONNAME}",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ParseManagedHSMPrivateEndpointConnectionID(v.Input)
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
		if actual.ManagedHSMName != v.Expected.ManagedHSMName {
			t.Fatalf("Expected %q but got %q for ManagedHSMName", v.Expected.ManagedHSMName, actual.ManagedHSMName)
		}
		if actual.PrivateEndpointConnectionName != v.Expected.PrivateEndpointConnectionName {
			t.Fatalf("Expected %q but got %q for PrivateEndpointConnectionName", v.Expected.PrivateEndpointConnectionName, actual.PrivateEndpointConnectionName)
		}
	}
}

func TestParseManagedHSMPrivateEndpointConnectionIDInsensitively(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *ManagedHSMPrivateEndpointConnectionId
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
			// missing ManagedHSMName
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.KeyVault/",
			Error: true,
		},

		{
			// missing value for ManagedHSMName
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.KeyVault/managedHSMs/",
			Error: true,
		},

		{
			// missing PrivateEndpointConnectionName
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.KeyVault/managedHSMs/{name}/",
			Error: true,
		},

		{
			// missing value for PrivateEndpointConnectionName
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.KeyVault/managedHSMs/{name}/privateEndpointConnections/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.KeyVault/managedHSMs/{name}/privateEndpointConnections/{privateEndpointConnectionName}",
			Expected: &ManagedHSMPrivateEndpointConnectionId{
				SubscriptionId:                "{subscriptionId}",
				ResourceGroup:                 "{resourceGroupName}",
				ManagedHSMName:                "{name}",
				PrivateEndpointConnectionName: "{privateEndpointConnectionName}",
			},
		},

		{
			// lower-cased segment names
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.KeyVault/managedhsms/{name}/privateendpointconnections/{privateEndpointConnectionName}",
			Expected: &ManagedHSMPrivateEndpointConnectionId{
				SubscriptionId:                "{subscriptionId}",
				ResourceGroup:                 "{resourceGroupName}",
				ManagedHSMName:                "{name}",
				PrivateEndpointConnectionName: "{privateEndpointConnectionName}",
			},
		},

		{
			// upper-cased segment names
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.KeyVault/MANAGEDHSMS/{name}/PRIVATEENDPOINTCONNECTIONS/{privateEndpointConnectionName}",
			Expected: &ManagedHSMPrivateEndpointConnectionId{
				SubscriptionId:                "{subscriptionId}",
				ResourceGroup:                 "{resourceGroupName}",
				ManagedHSMName:                "{name}",
				PrivateEndpointConnectionName: "{privateEndpointConnectionName}",
			},
		},

		{
			// mixed-cased segment names
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.KeyVault/MaNaGeDhSmS/{name}/PrIvAtEeNdPoInTcOnNeCtIoNs/{privateEndpointConnectionName}",
			Expected: &ManagedHSMPrivateEndpointConnectionId{
				SubscriptionId:                "{subscriptionId}",
				ResourceGroup:                 "{resourceGroupName}",
				ManagedHSMName:                "{name}",
				PrivateEndpointConnectionName: "{privateEndpointConnectionName}",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ParseManagedHSMPrivateEndpointConnectionIDInsensitively(v.Input)
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
		if actual.ManagedHSMName != v.Expected.ManagedHSMName {
			t.Fatalf("Expected %q but got %q for ManagedHSMName", v.Expected.ManagedHSMName, actual.ManagedHSMName)
		}
		if actual.PrivateEndpointConnectionName != v.Expected.PrivateEndpointConnectionName {
			t.Fatalf("Expected %q but got %q for PrivateEndpointConnectionName", v.Expected.PrivateEndpointConnectionName, actual.PrivateEndpointConnectionName)
		}
	}
}
