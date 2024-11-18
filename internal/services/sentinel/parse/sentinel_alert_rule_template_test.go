// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = SentinelAlertRuleTemplateId{}

func TestSentinelAlertRuleTemplateIDFormatter(t *testing.T) {
	actual := NewSentinelAlertRuleTemplateID("12345678-1234-9876-4563-123456789012", "resGroup1", "workspace1", "template1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1/providers/Microsoft.SecurityInsights/alertRuleTemplates/template1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestSentinelAlertRuleTemplateID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *SentinelAlertRuleTemplateId
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
			// missing WorkspaceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/",
			Error: true,
		},

		{
			// missing value for WorkspaceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/workspaces/",
			Error: true,
		},

		{
			// missing AlertRuleTemplateName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1/providers/Microsoft.SecurityInsights/",
			Error: true,
		},

		{
			// missing value for AlertRuleTemplateName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1/providers/Microsoft.SecurityInsights/alertRuleTemplates/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1/providers/Microsoft.SecurityInsights/alertRuleTemplates/template1",
			Expected: &SentinelAlertRuleTemplateId{
				SubscriptionId:        "12345678-1234-9876-4563-123456789012",
				ResourceGroup:         "resGroup1",
				WorkspaceName:         "workspace1",
				AlertRuleTemplateName: "template1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.OPERATIONALINSIGHTS/WORKSPACES/WORKSPACE1/PROVIDERS/MICROSOFT.SECURITYINSIGHTS/ALERTRULETEMPLATES/TEMPLATE1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := SentinelAlertRuleTemplateID(v.Input)
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
		if actual.WorkspaceName != v.Expected.WorkspaceName {
			t.Fatalf("Expected %q but got %q for WorkspaceName", v.Expected.WorkspaceName, actual.WorkspaceName)
		}
		if actual.AlertRuleTemplateName != v.Expected.AlertRuleTemplateName {
			t.Fatalf("Expected %q but got %q for AlertRuleTemplateName", v.Expected.AlertRuleTemplateName, actual.AlertRuleTemplateName)
		}
	}
}

func TestSentinelAlertRuleTemplateIDInsensitively(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *SentinelAlertRuleTemplateId
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
			// missing WorkspaceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/",
			Error: true,
		},

		{
			// missing value for WorkspaceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/workspaces/",
			Error: true,
		},

		{
			// missing AlertRuleTemplateName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1/providers/Microsoft.SecurityInsights/",
			Error: true,
		},

		{
			// missing value for AlertRuleTemplateName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1/providers/Microsoft.SecurityInsights/alertRuleTemplates/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1/providers/Microsoft.SecurityInsights/alertRuleTemplates/template1",
			Expected: &SentinelAlertRuleTemplateId{
				SubscriptionId:        "12345678-1234-9876-4563-123456789012",
				ResourceGroup:         "resGroup1",
				WorkspaceName:         "workspace1",
				AlertRuleTemplateName: "template1",
			},
		},

		{
			// lower-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1/providers/Microsoft.SecurityInsights/alertruletemplates/template1",
			Expected: &SentinelAlertRuleTemplateId{
				SubscriptionId:        "12345678-1234-9876-4563-123456789012",
				ResourceGroup:         "resGroup1",
				WorkspaceName:         "workspace1",
				AlertRuleTemplateName: "template1",
			},
		},

		{
			// upper-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/WORKSPACES/workspace1/providers/Microsoft.SecurityInsights/ALERTRULETEMPLATES/template1",
			Expected: &SentinelAlertRuleTemplateId{
				SubscriptionId:        "12345678-1234-9876-4563-123456789012",
				ResourceGroup:         "resGroup1",
				WorkspaceName:         "workspace1",
				AlertRuleTemplateName: "template1",
			},
		},

		{
			// mixed-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/WoRkSpAcEs/workspace1/providers/Microsoft.SecurityInsights/AlErTrUlEtEmPlAtEs/template1",
			Expected: &SentinelAlertRuleTemplateId{
				SubscriptionId:        "12345678-1234-9876-4563-123456789012",
				ResourceGroup:         "resGroup1",
				WorkspaceName:         "workspace1",
				AlertRuleTemplateName: "template1",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := SentinelAlertRuleTemplateIDInsensitively(v.Input)
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
		if actual.WorkspaceName != v.Expected.WorkspaceName {
			t.Fatalf("Expected %q but got %q for WorkspaceName", v.Expected.WorkspaceName, actual.WorkspaceName)
		}
		if actual.AlertRuleTemplateName != v.Expected.AlertRuleTemplateName {
			t.Fatalf("Expected %q but got %q for AlertRuleTemplateName", v.Expected.AlertRuleTemplateName, actual.AlertRuleTemplateName)
		}
	}
}
