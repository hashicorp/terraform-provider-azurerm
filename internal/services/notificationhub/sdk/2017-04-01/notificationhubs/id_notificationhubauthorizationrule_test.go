package notificationhubs

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = NotificationHubAuthorizationRuleId{}

func TestNewNotificationHubAuthorizationRuleID(t *testing.T) {
	id := NewNotificationHubAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceValue", "notificationHubValue", "authorizationRuleValue")

	if id.SubscriptionId != "12345678-1234-9876-4563-123456789012" {
		t.Fatalf("Expected %q but got %q for Segment 'SubscriptionId'", id.SubscriptionId, "12345678-1234-9876-4563-123456789012")
	}

	if id.ResourceGroupName != "example-resource-group" {
		t.Fatalf("Expected %q but got %q for Segment 'ResourceGroupName'", id.ResourceGroupName, "example-resource-group")
	}

	if id.NamespaceName != "namespaceValue" {
		t.Fatalf("Expected %q but got %q for Segment 'NamespaceName'", id.NamespaceName, "namespaceValue")
	}

	if id.NotificationHubName != "notificationHubValue" {
		t.Fatalf("Expected %q but got %q for Segment 'NotificationHubName'", id.NotificationHubName, "notificationHubValue")
	}

	if id.AuthorizationRuleName != "authorizationRuleValue" {
		t.Fatalf("Expected %q but got %q for Segment 'AuthorizationRuleName'", id.AuthorizationRuleName, "authorizationRuleValue")
	}
}

func TestFormatNotificationHubAuthorizationRuleID(t *testing.T) {
	actual := NewNotificationHubAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceValue", "notificationHubValue", "authorizationRuleValue").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.NotificationHubs/namespaces/namespaceValue/notificationHubs/notificationHubValue/authorizationRules/authorizationRuleValue"
	if actual != expected {
		t.Fatalf("Expected the Formatted ID to be %q but got %q", expected, actual)
	}
}

func TestParseNotificationHubAuthorizationRuleID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *NotificationHubAuthorizationRuleId
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
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.NotificationHubs",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.NotificationHubs/namespaces",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.NotificationHubs/namespaces/namespaceValue",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.NotificationHubs/namespaces/namespaceValue/notificationHubs",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.NotificationHubs/namespaces/namespaceValue/notificationHubs/notificationHubValue",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.NotificationHubs/namespaces/namespaceValue/notificationHubs/notificationHubValue/authorizationRules",
			Error: true,
		},
		{
			// Valid URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.NotificationHubs/namespaces/namespaceValue/notificationHubs/notificationHubValue/authorizationRules/authorizationRuleValue",
			Expected: &NotificationHubAuthorizationRuleId{
				SubscriptionId:        "12345678-1234-9876-4563-123456789012",
				ResourceGroupName:     "example-resource-group",
				NamespaceName:         "namespaceValue",
				NotificationHubName:   "notificationHubValue",
				AuthorizationRuleName: "authorizationRuleValue",
			},
		},
		{
			// Invalid (Valid Uri with Extra segment)
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.NotificationHubs/namespaces/namespaceValue/notificationHubs/notificationHubValue/authorizationRules/authorizationRuleValue/extra",
			Error: true,
		},
	}
	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ParseNotificationHubAuthorizationRuleID(v.Input)
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

		if actual.NotificationHubName != v.Expected.NotificationHubName {
			t.Fatalf("Expected %q but got %q for NotificationHubName", v.Expected.NotificationHubName, actual.NotificationHubName)
		}

		if actual.AuthorizationRuleName != v.Expected.AuthorizationRuleName {
			t.Fatalf("Expected %q but got %q for AuthorizationRuleName", v.Expected.AuthorizationRuleName, actual.AuthorizationRuleName)
		}

	}
}

func TestParseNotificationHubAuthorizationRuleIDInsensitively(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *NotificationHubAuthorizationRuleId
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
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.NotificationHubs",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/eXaMpLe-rEsOuRcE-GrOuP/pRoViDeRs/mIcRoSoFt.nOtIfIcAtIoNhUbS",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.NotificationHubs/namespaces",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/eXaMpLe-rEsOuRcE-GrOuP/pRoViDeRs/mIcRoSoFt.nOtIfIcAtIoNhUbS/nAmEsPaCeS",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.NotificationHubs/namespaces/namespaceValue",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/eXaMpLe-rEsOuRcE-GrOuP/pRoViDeRs/mIcRoSoFt.nOtIfIcAtIoNhUbS/nAmEsPaCeS/nAmEsPaCeVaLuE",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.NotificationHubs/namespaces/namespaceValue/notificationHubs",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/eXaMpLe-rEsOuRcE-GrOuP/pRoViDeRs/mIcRoSoFt.nOtIfIcAtIoNhUbS/nAmEsPaCeS/nAmEsPaCeVaLuE/nOtIfIcAtIoNhUbS",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.NotificationHubs/namespaces/namespaceValue/notificationHubs/notificationHubValue",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/eXaMpLe-rEsOuRcE-GrOuP/pRoViDeRs/mIcRoSoFt.nOtIfIcAtIoNhUbS/nAmEsPaCeS/nAmEsPaCeVaLuE/nOtIfIcAtIoNhUbS/nOtIfIcAtIoNhUbVaLuE",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.NotificationHubs/namespaces/namespaceValue/notificationHubs/notificationHubValue/authorizationRules",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/eXaMpLe-rEsOuRcE-GrOuP/pRoViDeRs/mIcRoSoFt.nOtIfIcAtIoNhUbS/nAmEsPaCeS/nAmEsPaCeVaLuE/nOtIfIcAtIoNhUbS/nOtIfIcAtIoNhUbVaLuE/aUtHoRiZaTiOnRuLeS",
			Error: true,
		},
		{
			// Valid URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.NotificationHubs/namespaces/namespaceValue/notificationHubs/notificationHubValue/authorizationRules/authorizationRuleValue",
			Expected: &NotificationHubAuthorizationRuleId{
				SubscriptionId:        "12345678-1234-9876-4563-123456789012",
				ResourceGroupName:     "example-resource-group",
				NamespaceName:         "namespaceValue",
				NotificationHubName:   "notificationHubValue",
				AuthorizationRuleName: "authorizationRuleValue",
			},
		},
		{
			// Invalid (Valid Uri with Extra segment)
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.NotificationHubs/namespaces/namespaceValue/notificationHubs/notificationHubValue/authorizationRules/authorizationRuleValue/extra",
			Error: true,
		},
		{
			// Valid URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/eXaMpLe-rEsOuRcE-GrOuP/pRoViDeRs/mIcRoSoFt.nOtIfIcAtIoNhUbS/nAmEsPaCeS/nAmEsPaCeVaLuE/nOtIfIcAtIoNhUbS/nOtIfIcAtIoNhUbVaLuE/aUtHoRiZaTiOnRuLeS/aUtHoRiZaTiOnRuLeVaLuE",
			Expected: &NotificationHubAuthorizationRuleId{
				SubscriptionId:        "12345678-1234-9876-4563-123456789012",
				ResourceGroupName:     "eXaMpLe-rEsOuRcE-GrOuP",
				NamespaceName:         "nAmEsPaCeVaLuE",
				NotificationHubName:   "nOtIfIcAtIoNhUbVaLuE",
				AuthorizationRuleName: "aUtHoRiZaTiOnRuLeVaLuE",
			},
		},
		{
			// Invalid (Valid Uri with Extra segment - mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/eXaMpLe-rEsOuRcE-GrOuP/pRoViDeRs/mIcRoSoFt.nOtIfIcAtIoNhUbS/nAmEsPaCeS/nAmEsPaCeVaLuE/nOtIfIcAtIoNhUbS/nOtIfIcAtIoNhUbVaLuE/aUtHoRiZaTiOnRuLeS/aUtHoRiZaTiOnRuLeVaLuE/extra",
			Error: true,
		},
	}
	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ParseNotificationHubAuthorizationRuleIDInsensitively(v.Input)
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

		if actual.NotificationHubName != v.Expected.NotificationHubName {
			t.Fatalf("Expected %q but got %q for NotificationHubName", v.Expected.NotificationHubName, actual.NotificationHubName)
		}

		if actual.AuthorizationRuleName != v.Expected.AuthorizationRuleName {
			t.Fatalf("Expected %q but got %q for AuthorizationRuleName", v.Expected.AuthorizationRuleName, actual.AuthorizationRuleName)
		}

	}
}

func TestSegmentsForNotificationHubAuthorizationRuleId(t *testing.T) {
	segments := NotificationHubAuthorizationRuleId{}.Segments()
	if len(segments) == 0 {
		t.Fatalf("NotificationHubAuthorizationRuleId has no segments")
	}

	uniqueNames := make(map[string]struct{}, 0)
	for _, segment := range segments {
		uniqueNames[segment.Name] = struct{}{}
	}
	if len(uniqueNames) != len(segments) {
		t.Fatalf("Expected the Segments to be unique but got %q unique segments and %d total segments", len(uniqueNames), len(segments))
	}
}
