package notificationrecipientemail

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &NotificationRecipientEmailId{}

func TestNewNotificationRecipientEmailID(t *testing.T) {
	id := NewNotificationRecipientEmailID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId", "AccountClosedPublisher", "recipientEmailName")

	if id.SubscriptionId != "12345678-1234-9876-4563-123456789012" {
		t.Fatalf("Expected %q but got %q for Segment 'SubscriptionId'", id.SubscriptionId, "12345678-1234-9876-4563-123456789012")
	}

	if id.ResourceGroupName != "example-resource-group" {
		t.Fatalf("Expected %q but got %q for Segment 'ResourceGroupName'", id.ResourceGroupName, "example-resource-group")
	}

	if id.ServiceName != "serviceName" {
		t.Fatalf("Expected %q but got %q for Segment 'ServiceName'", id.ServiceName, "serviceName")
	}

	if id.WorkspaceId != "workspaceId" {
		t.Fatalf("Expected %q but got %q for Segment 'WorkspaceId'", id.WorkspaceId, "workspaceId")
	}

	if id.NotificationName != "AccountClosedPublisher" {
		t.Fatalf("Expected %q but got %q for Segment 'NotificationName'", id.NotificationName, "AccountClosedPublisher")
	}

	if id.RecipientEmailName != "recipientEmailName" {
		t.Fatalf("Expected %q but got %q for Segment 'RecipientEmailName'", id.RecipientEmailName, "recipientEmailName")
	}
}

func TestFormatNotificationRecipientEmailID(t *testing.T) {
	actual := NewNotificationRecipientEmailID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId", "AccountClosedPublisher", "recipientEmailName").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.ApiManagement/service/serviceName/workspaces/workspaceId/notifications/AccountClosedPublisher/recipientEmails/recipientEmailName"
	if actual != expected {
		t.Fatalf("Expected the Formatted ID to be %q but got %q", expected, actual)
	}
}

func TestParseNotificationRecipientEmailID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *NotificationRecipientEmailId
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
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.ApiManagement",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.ApiManagement/service",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.ApiManagement/service/serviceName",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.ApiManagement/service/serviceName/workspaces",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.ApiManagement/service/serviceName/workspaces/workspaceId",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.ApiManagement/service/serviceName/workspaces/workspaceId/notifications",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.ApiManagement/service/serviceName/workspaces/workspaceId/notifications/AccountClosedPublisher",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.ApiManagement/service/serviceName/workspaces/workspaceId/notifications/AccountClosedPublisher/recipientEmails",
			Error: true,
		},
		{
			// Valid URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.ApiManagement/service/serviceName/workspaces/workspaceId/notifications/AccountClosedPublisher/recipientEmails/recipientEmailName",
			Expected: &NotificationRecipientEmailId{
				SubscriptionId:     "12345678-1234-9876-4563-123456789012",
				ResourceGroupName:  "example-resource-group",
				ServiceName:        "serviceName",
				WorkspaceId:        "workspaceId",
				NotificationName:   "AccountClosedPublisher",
				RecipientEmailName: "recipientEmailName",
			},
		},
		{
			// Invalid (Valid Uri with Extra segment)
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.ApiManagement/service/serviceName/workspaces/workspaceId/notifications/AccountClosedPublisher/recipientEmails/recipientEmailName/extra",
			Error: true,
		},
	}
	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ParseNotificationRecipientEmailID(v.Input)
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

		if actual.ServiceName != v.Expected.ServiceName {
			t.Fatalf("Expected %q but got %q for ServiceName", v.Expected.ServiceName, actual.ServiceName)
		}

		if actual.WorkspaceId != v.Expected.WorkspaceId {
			t.Fatalf("Expected %q but got %q for WorkspaceId", v.Expected.WorkspaceId, actual.WorkspaceId)
		}

		if actual.NotificationName != v.Expected.NotificationName {
			t.Fatalf("Expected %q but got %q for NotificationName", v.Expected.NotificationName, actual.NotificationName)
		}

		if actual.RecipientEmailName != v.Expected.RecipientEmailName {
			t.Fatalf("Expected %q but got %q for RecipientEmailName", v.Expected.RecipientEmailName, actual.RecipientEmailName)
		}

	}
}

func TestParseNotificationRecipientEmailIDInsensitively(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *NotificationRecipientEmailId
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
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.ApiManagement",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/eXaMpLe-rEsOuRcE-GrOuP/pRoViDeRs/mIcRoSoFt.aPiMaNaGeMeNt",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.ApiManagement/service",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/eXaMpLe-rEsOuRcE-GrOuP/pRoViDeRs/mIcRoSoFt.aPiMaNaGeMeNt/sErViCe",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.ApiManagement/service/serviceName",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/eXaMpLe-rEsOuRcE-GrOuP/pRoViDeRs/mIcRoSoFt.aPiMaNaGeMeNt/sErViCe/sErViCeNaMe",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.ApiManagement/service/serviceName/workspaces",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/eXaMpLe-rEsOuRcE-GrOuP/pRoViDeRs/mIcRoSoFt.aPiMaNaGeMeNt/sErViCe/sErViCeNaMe/wOrKsPaCeS",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.ApiManagement/service/serviceName/workspaces/workspaceId",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/eXaMpLe-rEsOuRcE-GrOuP/pRoViDeRs/mIcRoSoFt.aPiMaNaGeMeNt/sErViCe/sErViCeNaMe/wOrKsPaCeS/wOrKsPaCeId",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.ApiManagement/service/serviceName/workspaces/workspaceId/notifications",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/eXaMpLe-rEsOuRcE-GrOuP/pRoViDeRs/mIcRoSoFt.aPiMaNaGeMeNt/sErViCe/sErViCeNaMe/wOrKsPaCeS/wOrKsPaCeId/nOtIfIcAtIoNs",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.ApiManagement/service/serviceName/workspaces/workspaceId/notifications/AccountClosedPublisher",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/eXaMpLe-rEsOuRcE-GrOuP/pRoViDeRs/mIcRoSoFt.aPiMaNaGeMeNt/sErViCe/sErViCeNaMe/wOrKsPaCeS/wOrKsPaCeId/nOtIfIcAtIoNs/aCcOuNtClOsEdPuBlIsHeR",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.ApiManagement/service/serviceName/workspaces/workspaceId/notifications/AccountClosedPublisher/recipientEmails",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/eXaMpLe-rEsOuRcE-GrOuP/pRoViDeRs/mIcRoSoFt.aPiMaNaGeMeNt/sErViCe/sErViCeNaMe/wOrKsPaCeS/wOrKsPaCeId/nOtIfIcAtIoNs/aCcOuNtClOsEdPuBlIsHeR/rEcIpIeNtEmAiLs",
			Error: true,
		},
		{
			// Valid URI
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.ApiManagement/service/serviceName/workspaces/workspaceId/notifications/AccountClosedPublisher/recipientEmails/recipientEmailName",
			Expected: &NotificationRecipientEmailId{
				SubscriptionId:     "12345678-1234-9876-4563-123456789012",
				ResourceGroupName:  "example-resource-group",
				ServiceName:        "serviceName",
				WorkspaceId:        "workspaceId",
				NotificationName:   "AccountClosedPublisher",
				RecipientEmailName: "recipientEmailName",
			},
		},
		{
			// Invalid (Valid Uri with Extra segment)
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.ApiManagement/service/serviceName/workspaces/workspaceId/notifications/AccountClosedPublisher/recipientEmails/recipientEmailName/extra",
			Error: true,
		},
		{
			// Valid URI (mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/eXaMpLe-rEsOuRcE-GrOuP/pRoViDeRs/mIcRoSoFt.aPiMaNaGeMeNt/sErViCe/sErViCeNaMe/wOrKsPaCeS/wOrKsPaCeId/nOtIfIcAtIoNs/aCcOuNtClOsEdPuBlIsHeR/rEcIpIeNtEmAiLs/rEcIpIeNtEmAiLnAmE",
			Expected: &NotificationRecipientEmailId{
				SubscriptionId:     "12345678-1234-9876-4563-123456789012",
				ResourceGroupName:  "eXaMpLe-rEsOuRcE-GrOuP",
				ServiceName:        "sErViCeNaMe",
				WorkspaceId:        "wOrKsPaCeId",
				NotificationName:   "AccountClosedPublisher",
				RecipientEmailName: "rEcIpIeNtEmAiLnAmE",
			},
		},
		{
			// Invalid (Valid Uri with Extra segment - mIxEd CaSe since this is insensitive)
			Input: "/sUbScRiPtIoNs/12345678-1234-9876-4563-123456789012/rEsOuRcEgRoUpS/eXaMpLe-rEsOuRcE-GrOuP/pRoViDeRs/mIcRoSoFt.aPiMaNaGeMeNt/sErViCe/sErViCeNaMe/wOrKsPaCeS/wOrKsPaCeId/nOtIfIcAtIoNs/aCcOuNtClOsEdPuBlIsHeR/rEcIpIeNtEmAiLs/rEcIpIeNtEmAiLnAmE/extra",
			Error: true,
		},
	}
	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ParseNotificationRecipientEmailIDInsensitively(v.Input)
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

		if actual.ServiceName != v.Expected.ServiceName {
			t.Fatalf("Expected %q but got %q for ServiceName", v.Expected.ServiceName, actual.ServiceName)
		}

		if actual.WorkspaceId != v.Expected.WorkspaceId {
			t.Fatalf("Expected %q but got %q for WorkspaceId", v.Expected.WorkspaceId, actual.WorkspaceId)
		}

		if actual.NotificationName != v.Expected.NotificationName {
			t.Fatalf("Expected %q but got %q for NotificationName", v.Expected.NotificationName, actual.NotificationName)
		}

		if actual.RecipientEmailName != v.Expected.RecipientEmailName {
			t.Fatalf("Expected %q but got %q for RecipientEmailName", v.Expected.RecipientEmailName, actual.RecipientEmailName)
		}

	}
}

func TestSegmentsForNotificationRecipientEmailId(t *testing.T) {
	segments := NotificationRecipientEmailId{}.Segments()
	if len(segments) == 0 {
		t.Fatalf("NotificationRecipientEmailId has no segments")
	}

	uniqueNames := make(map[string]struct{}, 0)
	for _, segment := range segments {
		uniqueNames[segment.Name] = struct{}{}
	}
	if len(uniqueNames) != len(segments) {
		t.Fatalf("Expected the Segments to be unique but got %q unique segments and %d total segments", len(uniqueNames), len(segments))
	}
}
