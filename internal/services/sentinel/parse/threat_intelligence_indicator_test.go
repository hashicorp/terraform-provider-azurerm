// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = ThreatIntelligenceIndicatorId{}

func TestThreatIntelligenceIndicatorIDFormatter(t *testing.T) {
	actual := NewThreatIntelligenceIndicatorID("12345678-1234-9876-4563-123456789012", "resGroup1", "workspace1", "main", "indicator1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1/providers/Microsoft.SecurityInsights/threatIntelligence/main/indicators/indicator1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestThreatIntelligenceIndicatorID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *ThreatIntelligenceIndicatorId
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
			// missing ThreatIntelligenceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1/providers/Microsoft.SecurityInsights/",
			Error: true,
		},

		{
			// missing value for ThreatIntelligenceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1/providers/Microsoft.SecurityInsights/threatIntelligence/",
			Error: true,
		},

		{
			// missing IndicatorName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1/providers/Microsoft.SecurityInsights/threatIntelligence/main/",
			Error: true,
		},

		{
			// missing value for IndicatorName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1/providers/Microsoft.SecurityInsights/threatIntelligence/main/indicators/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1/providers/Microsoft.SecurityInsights/threatIntelligence/main/indicators/indicator1",
			Expected: &ThreatIntelligenceIndicatorId{
				SubscriptionId:         "12345678-1234-9876-4563-123456789012",
				ResourceGroup:          "resGroup1",
				WorkspaceName:          "workspace1",
				ThreatIntelligenceName: "main",
				IndicatorName:          "indicator1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.OPERATIONALINSIGHTS/WORKSPACES/WORKSPACE1/PROVIDERS/MICROSOFT.SECURITYINSIGHTS/THREATINTELLIGENCE/MAIN/INDICATORS/INDICATOR1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ThreatIntelligenceIndicatorID(v.Input)
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
		if actual.ThreatIntelligenceName != v.Expected.ThreatIntelligenceName {
			t.Fatalf("Expected %q but got %q for ThreatIntelligenceName", v.Expected.ThreatIntelligenceName, actual.ThreatIntelligenceName)
		}
		if actual.IndicatorName != v.Expected.IndicatorName {
			t.Fatalf("Expected %q but got %q for IndicatorName", v.Expected.IndicatorName, actual.IndicatorName)
		}
	}
}
