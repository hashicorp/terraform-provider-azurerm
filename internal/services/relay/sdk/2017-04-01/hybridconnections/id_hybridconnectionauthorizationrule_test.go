package hybridconnections

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = HybridConnectionAuthorizationRuleId{}

func TestNewHybridConnectionAuthorizationRuleID(t *testing.T) {
	id := NewHybridConnectionAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceValue", "hybridConnectionValue", "authorizationRuleValue")

	if id.SubscriptionId != "12345678-1234-9876-4563-123456789012" {
		t.Fatalf("Expected %q but got %q for Segment 'SubscriptionId'", id.SubscriptionId, "12345678-1234-9876-4563-123456789012")
	}

	if id.ResourceGroupName != "example-resource-group" {
		t.Fatalf("Expected %q but got %q for Segment 'ResourceGroupName'", id.ResourceGroupName, "example-resource-group")
	}

	if id.NamespaceName != "namespaceValue" {
		t.Fatalf("Expected %q but got %q for Segment 'NamespaceName'", id.NamespaceName, "namespaceValue")
	}

	if id.HybridConnectionName != "hybridConnectionValue" {
		t.Fatalf("Expected %q but got %q for Segment 'HybridConnectionName'", id.HybridConnectionName, "hybridConnectionValue")
	}

	if id.AuthorizationRuleName != "authorizationRuleValue" {
		t.Fatalf("Expected %q but got %q for Segment 'AuthorizationRuleName'", id.AuthorizationRuleName, "authorizationRuleValue")
	}
}

func TestFormatHybridConnectionAuthorizationRuleID(t *testing.T) {
	actual := NewHybridConnectionAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceValue", "hybridConnectionValue", "authorizationRuleValue").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.Relay/namespaces/namespaceValue/hybridConnections/hybridConnectionValue/authorizationRules/authorizationRuleValue"
	if actual != expected {
		t.Fatalf("Expected the Formatted ID to be %q but got %q", actual, expected)
	}
}

func TestParseHybridConnectionAuthorizationRuleID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *HybridConnectionAuthorizationRuleId
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
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.Relay",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.Relay/namespaces",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.Relay/namespaces/namespaceValue",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.Relay/namespaces/namespaceValue/hybridConnections",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.Relay/namespaces/namespaceValue/hybridConnections/hybridConnectionValue",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.Relay/namespaces/namespaceValue/hybridConnections/hybridConnectionValue/authorizationRules",
			Error: true,
		},
		{
			// Valid URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.Relay/namespaces/namespaceValue/hybridConnections/hybridConnectionValue/authorizationRules/authorizationRuleValue",
			Expected: &HybridConnectionAuthorizationRuleId{
				SubscriptionId:        "12345678-1234-9876-4563-123456789012",
				ResourceGroupName:     "example-resource-group",
				NamespaceName:         "namespaceValue",
				HybridConnectionName:  "hybridConnectionValue",
				AuthorizationRuleName: "authorizationRuleValue",
			},
		},
		{
			// Invalid (Valid Uri with Extra segment)
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.Relay/namespaces/namespaceValue/hybridConnections/hybridConnectionValue/authorizationRules/authorizationRuleValue/extra",
			Error: true,
		},
	}
	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ParseHybridConnectionAuthorizationRuleID(v.Input)
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

		if actual.NamespaceName != v.Expected.NamespaceName {
			t.Fatalf("Expected %q but got %q for NamespaceName", v.Expected.NamespaceName, actual.NamespaceName)
		}

		if actual.HybridConnectionName != v.Expected.HybridConnectionName {
			t.Fatalf("Expected %q but got %q for HybridConnectionName", v.Expected.HybridConnectionName, actual.HybridConnectionName)
		}

		if actual.AuthorizationRuleName != v.Expected.AuthorizationRuleName {
			t.Fatalf("Expected %q but got %q for AuthorizationRuleName", v.Expected.AuthorizationRuleName, actual.AuthorizationRuleName)
		}

	}
}

func TestParseHybridConnectionAuthorizationRuleIDInsensitively(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *HybridConnectionAuthorizationRuleId
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
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.Relay",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/eXaMpLe-rEsOuRcE-GrOuP/pRoViDeRs/mIcRoSoFt.rElAy",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.Relay/namespaces",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/eXaMpLe-rEsOuRcE-GrOuP/pRoViDeRs/mIcRoSoFt.rElAy/nAmEsPaCeS",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.Relay/namespaces/namespaceValue",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/eXaMpLe-rEsOuRcE-GrOuP/pRoViDeRs/mIcRoSoFt.rElAy/nAmEsPaCeS/nAmEsPaCeVaLuE",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.Relay/namespaces/namespaceValue/hybridConnections",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/eXaMpLe-rEsOuRcE-GrOuP/pRoViDeRs/mIcRoSoFt.rElAy/nAmEsPaCeS/nAmEsPaCeVaLuE/hYbRiDcOnNeCtIoNs",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.Relay/namespaces/namespaceValue/hybridConnections/hybridConnectionValue",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/eXaMpLe-rEsOuRcE-GrOuP/pRoViDeRs/mIcRoSoFt.rElAy/nAmEsPaCeS/nAmEsPaCeVaLuE/hYbRiDcOnNeCtIoNs/hYbRiDcOnNeCtIoNvAlUe",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.Relay/namespaces/namespaceValue/hybridConnections/hybridConnectionValue/authorizationRules",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/eXaMpLe-rEsOuRcE-GrOuP/pRoViDeRs/mIcRoSoFt.rElAy/nAmEsPaCeS/nAmEsPaCeVaLuE/hYbRiDcOnNeCtIoNs/hYbRiDcOnNeCtIoNvAlUe/aUtHoRiZaTiOnRuLeS",
			Error: true,
		},
		{
			// Valid URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.Relay/namespaces/namespaceValue/hybridConnections/hybridConnectionValue/authorizationRules/authorizationRuleValue",
			Expected: &HybridConnectionAuthorizationRuleId{
				SubscriptionId:        "12345678-1234-9876-4563-123456789012",
				ResourceGroupName:     "example-resource-group",
				NamespaceName:         "namespaceValue",
				HybridConnectionName:  "hybridConnectionValue",
				AuthorizationRuleName: "authorizationRuleValue",
			},
		},
		{
			// Invalid (Valid Uri with Extra segment)
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.Relay/namespaces/namespaceValue/hybridConnections/hybridConnectionValue/authorizationRules/authorizationRuleValue/extra",
			Error: true,
		},
		{
			// Valid URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/eXaMpLe-rEsOuRcE-GrOuP/pRoViDeRs/mIcRoSoFt.rElAy/nAmEsPaCeS/nAmEsPaCeVaLuE/hYbRiDcOnNeCtIoNs/hYbRiDcOnNeCtIoNvAlUe/aUtHoRiZaTiOnRuLeS/aUtHoRiZaTiOnRuLeVaLuE",
			Expected: &HybridConnectionAuthorizationRuleId{
				SubscriptionId:        "12345678-1234-9876-4563-123456789012",
				ResourceGroupName:     "eXaMpLe-rEsOuRcE-GrOuP",
				NamespaceName:         "nAmEsPaCeVaLuE",
				HybridConnectionName:  "hYbRiDcOnNeCtIoNvAlUe",
				AuthorizationRuleName: "aUtHoRiZaTiOnRuLeVaLuE",
			},
		},
		{
			// Invalid (Valid Uri with Extra segment - mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/eXaMpLe-rEsOuRcE-GrOuP/pRoViDeRs/mIcRoSoFt.rElAy/nAmEsPaCeS/nAmEsPaCeVaLuE/hYbRiDcOnNeCtIoNs/hYbRiDcOnNeCtIoNvAlUe/aUtHoRiZaTiOnRuLeS/aUtHoRiZaTiOnRuLeVaLuE/extra",
			Error: true,
		},
	}
	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ParseHybridConnectionAuthorizationRuleIDInsensitively(v.Input)
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

		if actual.NamespaceName != v.Expected.NamespaceName {
			t.Fatalf("Expected %q but got %q for NamespaceName", v.Expected.NamespaceName, actual.NamespaceName)
		}

		if actual.HybridConnectionName != v.Expected.HybridConnectionName {
			t.Fatalf("Expected %q but got %q for HybridConnectionName", v.Expected.HybridConnectionName, actual.HybridConnectionName)
		}

		if actual.AuthorizationRuleName != v.Expected.AuthorizationRuleName {
			t.Fatalf("Expected %q but got %q for AuthorizationRuleName", v.Expected.AuthorizationRuleName, actual.AuthorizationRuleName)
		}

	}
}
