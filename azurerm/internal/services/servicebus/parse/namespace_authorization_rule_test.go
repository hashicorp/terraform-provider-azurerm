package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = NamespaceAuthorizationRuleId{}

func TestNamespaceAuthorizationRuleIDFormatter(t *testing.T) {
	actual := NewNamespaceAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "resGroup1", "namespace1", "authorizationRule1").ID("")
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ServiceBus/namespaces/namespace1/AuthorizationRules/authorizationRule1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestNamespaceAuthorizationRuleID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *NamespaceAuthorizationRuleId
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
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/",
			Error: true,
		},

		{
			// missing value for ResourceGroup
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/",
			Error: true,
		},

		{
			// missing NamespaceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ServiceBus/",
			Error: true,
		},

		{
			// missing value for NamespaceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ServiceBus/namespaces/",
			Error: true,
		},

		{
			// missing AuthorizationRuleName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ServiceBus/namespaces/namespace1/",
			Error: true,
		},

		{
			// missing value for AuthorizationRuleName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ServiceBus/namespaces/namespace1/AuthorizationRules/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ServiceBus/namespaces/namespace1/AuthorizationRules/authorizationRule1",
			Expected: &NamespaceAuthorizationRuleId{
				SubscriptionId:        "12345678-1234-9876-4563-123456789012",
				ResourceGroup:         "resGroup1",
				NamespaceName:         "namespace1",
				AuthorizationRuleName: "authorizationRule1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.SERVICEBUS/NAMESPACES/NAMESPACE1/AUTHORIZATIONRULES/AUTHORIZATIONRULE1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := NamespaceAuthorizationRuleID(v.Input)
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
		if actual.NamespaceName != v.Expected.NamespaceName {
			t.Fatalf("Expected %q but got %q for NamespaceName", v.Expected.NamespaceName, actual.NamespaceName)
		}
		if actual.AuthorizationRuleName != v.Expected.AuthorizationRuleName {
			t.Fatalf("Expected %q but got %q for AuthorizationRuleName", v.Expected.AuthorizationRuleName, actual.AuthorizationRuleName)
		}
	}
}
