package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = SubscriptionRuleId{}

func TestSubscriptionRuleIDFormatter(t *testing.T) {
	actual := NewSubscriptionRuleID("12345678-1234-9876-4563-123456789012", "resGroup1", "namespace1", "topic1", "subscription1", "rule1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ServiceBus/namespaces/namespace1/topics/topic1/subscriptions/subscription1/rules/rule1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestSubscriptionRuleID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *SubscriptionRuleId
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
			// missing TopicName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ServiceBus/namespaces/namespace1/",
			Error: true,
		},

		{
			// missing value for TopicName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ServiceBus/namespaces/namespace1/topics/",
			Error: true,
		},

		{
			// missing SubscriptionName
			Input: "/",
			Error: true,
		},

		{
			// missing value for SubscriptionName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ServiceBus/namespaces/namespace1/topics/topic1/subscriptions/",
			Error: true,
		},

		{
			// missing RuleName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ServiceBus/namespaces/namespace1/topics/topic1/subscriptions/subscription1/",
			Error: true,
		},

		{
			// missing value for RuleName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ServiceBus/namespaces/namespace1/topics/topic1/subscriptions/subscription1/rules/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ServiceBus/namespaces/namespace1/topics/topic1/subscriptions/subscription1/rules/rule1",
			Expected: &SubscriptionRuleId{
				SubscriptionId:   "12345678-1234-9876-4563-123456789012",
				ResourceGroup:    "resGroup1",
				NamespaceName:    "namespace1",
				TopicName:        "topic1",
				SubscriptionName: "subscription1",
				RuleName:         "rule1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.SERVICEBUS/NAMESPACES/NAMESPACE1/TOPICS/TOPIC1/SUBSCRIPTIONS/SUBSCRIPTION1/RULES/RULE1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := SubscriptionRuleID(v.Input)
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
		if actual.TopicName != v.Expected.TopicName {
			t.Fatalf("Expected %q but got %q for TopicName", v.Expected.TopicName, actual.TopicName)
		}
		if actual.SubscriptionName != v.Expected.SubscriptionName {
			t.Fatalf("Expected %q but got %q for SubscriptionName", v.Expected.SubscriptionName, actual.SubscriptionName)
		}
		if actual.RuleName != v.Expected.RuleName {
			t.Fatalf("Expected %q but got %q for RuleName", v.Expected.RuleName, actual.RuleName)
		}
	}
}
