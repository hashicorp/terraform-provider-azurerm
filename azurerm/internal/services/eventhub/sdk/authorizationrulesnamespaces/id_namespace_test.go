package authorizationrulesnamespaces

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = NamespaceId{}

func TestNamespaceIDFormatter(t *testing.T) {
	actual := NewNamespaceID("{subscriptionId}", "{resourceGroupName}", "{namespaceName}").ID()
	expected := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventHub/namespaces/{namespaceName}"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestNamespaceID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *NamespaceId
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
			// missing Name
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventHub/",
			Error: true,
		},

		{
			// missing value for Name
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventHub/namespaces/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventHub/namespaces/{namespaceName}",
			Expected: &NamespaceId{
				SubscriptionId: "{subscriptionId}",
				ResourceGroup:  "{resourceGroupName}",
				Name:           "{namespaceName}",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/{SUBSCRIPTIONID}/RESOURCEGROUPS/{RESOURCEGROUPNAME}/PROVIDERS/MICROSOFT.EVENTHUB/NAMESPACES/{NAMESPACENAME}",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := NamespaceID(v.Input)
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
		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}
	}
}

func TestNamespaceIDInsensitively(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *NamespaceId
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
			// missing Name
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventHub/",
			Error: true,
		},

		{
			// missing value for Name
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventHub/namespaces/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventHub/namespaces/{namespaceName}",
			Expected: &NamespaceId{
				SubscriptionId: "{subscriptionId}",
				ResourceGroup:  "{resourceGroupName}",
				Name:           "{namespaceName}",
			},
		},

		{
			// lower-cased segment names
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventHub/namespaces/{namespaceName}",
			Expected: &NamespaceId{
				SubscriptionId: "{subscriptionId}",
				ResourceGroup:  "{resourceGroupName}",
				Name:           "{namespaceName}",
			},
		},

		{
			// upper-cased segment names
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventHub/NAMESPACES/{namespaceName}",
			Expected: &NamespaceId{
				SubscriptionId: "{subscriptionId}",
				ResourceGroup:  "{resourceGroupName}",
				Name:           "{namespaceName}",
			},
		},

		{
			// mixed-cased segment names
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventHub/NaMeSpAcEs/{namespaceName}",
			Expected: &NamespaceId{
				SubscriptionId: "{subscriptionId}",
				ResourceGroup:  "{resourceGroupName}",
				Name:           "{namespaceName}",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := NamespaceIDInsensitively(v.Input)
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
		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}
	}
}
