package parse

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = LogAnalyticsWorkspaceId{}

func TestSentinelAlertRuleIDFormatter(t *testing.T) {
	subscriptionId := "12345678-1234-5678-1234-123456789012"
	actual := NewLogAnalyticsWorkspaceID("group1", "space1").ID(subscriptionId)
	expected := "/subscriptions/12345678-1234-5678-1234-123456789012/resourcegroups/group1/providers/microsoft.operationalinsights/workspaces/space1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestLogAnalyticsWorkspaceID(t *testing.T) {
	testData := []struct {
		Name   string
		Input  string
		Error  bool
		Expect *LogAnalyticsWorkspaceId
	}{
		{
			Name:  "Empty",
			Input: "",
			Error: true,
		},
		{
			Name:  "No Resource Groups Segment",
			Input: "/subscriptions/11111111-1111-1111-1111-1111111111111",
			Error: true,
		},
		{
			Name:  "No Resource Groups Value",
			Input: "/subscriptions/11111111-1111-1111-1111-1111111111111/resourceGroups/",
			Error: true,
		},
		{
			Name:  "Resource Group ID",
			Input: "/subscriptions/11111111-1111-1111-1111-1111111111111/resourceGroups/resGroup1/",
			Error: true,
		},
		{
			Name:  "Missing Workspace Value",
			Input: "/subscriptions/11111111-1111-1111-1111-1111111111111/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/workspaces",
			Error: true,
		},
		{
			Name:  "Workspace Value",
			Input: "/subscriptions/11111111-1111-1111-1111-1111111111111/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1",
			Error: true,
			Expect: &LogAnalyticsWorkspaceId{
				ResourceGroup: "resGroup1",
				Name:          "workspace1",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := LogAnalyticsWorkspaceID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.ResourceGroup != v.Expect.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expect.ResourceGroup, actual.ResourceGroup)
		}

		if actual.Name != v.Expect.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expect.Name, actual.Name)
		}
	}
}
