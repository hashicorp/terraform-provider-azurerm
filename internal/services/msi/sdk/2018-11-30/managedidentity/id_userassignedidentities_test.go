package managedidentity

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = UserAssignedIdentitiesId{}

func TestUserAssignedIdentitiesIDFormatter(t *testing.T) {
	actual := NewUserAssignedIdentitiesID("{subscriptionId}", "{resourceGroupName}", "{resourceName}").ID()
	expected := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ManagedIdentity/userAssignedIdentities/{resourceName}"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestParseUserAssignedIdentitiesID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *UserAssignedIdentitiesId
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
			// missing UserAssignedIdentityName
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ManagedIdentity/",
			Error: true,
		},

		{
			// missing value for UserAssignedIdentityName
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ManagedIdentity/userAssignedIdentities/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ManagedIdentity/userAssignedIdentities/{resourceName}",
			Expected: &UserAssignedIdentitiesId{
				SubscriptionId:           "{subscriptionId}",
				ResourceGroup:            "{resourceGroupName}",
				UserAssignedIdentityName: "{resourceName}",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/{SUBSCRIPTIONID}/RESOURCEGROUPS/{RESOURCEGROUPNAME}/PROVIDERS/MICROSOFT.MANAGEDIDENTITY/USERASSIGNEDIDENTITIES/{RESOURCENAME}",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ParseUserAssignedIdentitiesID(v.Input)
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
		if actual.UserAssignedIdentityName != v.Expected.UserAssignedIdentityName {
			t.Fatalf("Expected %q but got %q for UserAssignedIdentityName", v.Expected.UserAssignedIdentityName, actual.UserAssignedIdentityName)
		}
	}
}

func TestParseUserAssignedIdentitiesIDInsensitively(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *UserAssignedIdentitiesId
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
			// missing UserAssignedIdentityName
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ManagedIdentity/",
			Error: true,
		},

		{
			// missing value for UserAssignedIdentityName
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ManagedIdentity/userAssignedIdentities/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ManagedIdentity/userAssignedIdentities/{resourceName}",
			Expected: &UserAssignedIdentitiesId{
				SubscriptionId:           "{subscriptionId}",
				ResourceGroup:            "{resourceGroupName}",
				UserAssignedIdentityName: "{resourceName}",
			},
		},

		{
			// lower-cased segment names
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ManagedIdentity/userassignedidentities/{resourceName}",
			Expected: &UserAssignedIdentitiesId{
				SubscriptionId:           "{subscriptionId}",
				ResourceGroup:            "{resourceGroupName}",
				UserAssignedIdentityName: "{resourceName}",
			},
		},

		{
			// upper-cased segment names
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ManagedIdentity/USERASSIGNEDIDENTITIES/{resourceName}",
			Expected: &UserAssignedIdentitiesId{
				SubscriptionId:           "{subscriptionId}",
				ResourceGroup:            "{resourceGroupName}",
				UserAssignedIdentityName: "{resourceName}",
			},
		},

		{
			// mixed-cased segment names
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ManagedIdentity/UsErAsSiGnEdIdEnTiTiEs/{resourceName}",
			Expected: &UserAssignedIdentitiesId{
				SubscriptionId:           "{subscriptionId}",
				ResourceGroup:            "{resourceGroupName}",
				UserAssignedIdentityName: "{resourceName}",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ParseUserAssignedIdentitiesIDInsensitively(v.Input)
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
		if actual.UserAssignedIdentityName != v.Expected.UserAssignedIdentityName {
			t.Fatalf("Expected %q but got %q for UserAssignedIdentityName", v.Expected.UserAssignedIdentityName, actual.UserAssignedIdentityName)
		}
	}
}
