package parse

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = VirtualDesktopWorkspaceApplicationGroupAssociationId{}

func TestVirtualDesktopWorkspaceApplicationGroupAssociationIDFormatter(t *testing.T) {
	subscriptionId := "12345678-1234-5678-1234-123456789012"
	workspaceId := NewVirtualDesktopWorkspaceId(subscriptionId, "group1", "workspace1")
	applicationGroupId := NewVirtualDesktopApplicationGroupId(subscriptionId, "group1", "appGroup1")

	actual := NewVirtualDesktopWorkspaceApplicationGroupAssociationId(workspaceId, applicationGroupId).ID("")
	expected := "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.DesktopVirtualization/workspaces/workspace1|/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.DesktopVirtualization/applicationgroups/appGroup1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestVirtualDesktopWorkspaceApplicationGroupAssociationID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *VirtualDesktopWorkspaceApplicationGroupAssociationId
	}{
		{
			Name:     "Empty",
			Input:    "",
			Expected: nil,
		},
		{
			Name:     "No Resource Groups Segment",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000",
			Expected: nil,
		},
		{
			Name:     "No Resource Groups Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/",
			Expected: nil,
		},
		{
			Name:     "Resource Group ID",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/",
			Expected: nil,
		},
		{
			Name:     "Missing Workspace Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/Microsoft.DesktopVirtualization/workspaces/",
			Expected: nil,
		},
		{
			Name:     "Virtual Desktop Workspace ID",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.DesktopVirtualization/workspaces/workspace1",
			Expected: nil,
		},
		{
			Name:     "Missing Application Group Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/Microsoft.DesktopVirtualization/applicationgroups/",
			Expected: nil,
		},
		{
			Name:     "Virtual Desktop Application Group ID",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.DesktopVirtualization/applicationgroups/appGroup1",
			Expected: nil,
		},
		{
			Name:  "Virtual Desktop Workspace ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.DesktopVirtualization/workspaces/workspace1|/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.DesktopVirtualization/applicationgroups/appGroup1",
			Expected: &VirtualDesktopWorkspaceApplicationGroupAssociationId{
				Workspace: VirtualDesktopWorkspaceId{
					SubscriptionId: "00000000-0000-0000-0000-000000000000",
					ResourceGroup:  "resGroup1",
					Name:           "workspace1",
				},
				ApplicationGroup: VirtualDesktopApplicationGroupId{
					SubscriptionId: "00000000-0000-0000-0000-000000000000",
					ResourceGroup:  "resGroup1",
					Name:           "appGroup1",
				},
			},
		},
		{
			Name:     "Wrong Casing",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.DesktopVirtualization/Workspaces/workspace1|/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.DesktopVirtualization/Applicationgroups/appGroup1",
			Expected: nil,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.Name)

		actual, err := VirtualDesktopWorkspaceApplicationGroupAssociationID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}
			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.ApplicationGroup.SubscriptionId != v.Expected.ApplicationGroup.SubscriptionId {
			t.Fatalf("Expected %q but got %q for ApplicationGroup.SubscriptionId", v.Expected.ApplicationGroup.SubscriptionId, actual.ApplicationGroup.SubscriptionId)
		}

		if actual.ApplicationGroup.ResourceGroup != v.Expected.ApplicationGroup.ResourceGroup {
			t.Fatalf("Expected %q but got %q for ApplicationGroup.ResourceGroup", v.Expected.ApplicationGroup.ResourceGroup, actual.ApplicationGroup.ResourceGroup)
		}

		if actual.ApplicationGroup.Name != v.Expected.ApplicationGroup.Name {
			t.Fatalf("Expected %q but got %q for ApplicationGroup.Name", v.Expected.ApplicationGroup.Name, actual.ApplicationGroup.Name)
		}

		if actual.Workspace.SubscriptionId != v.Expected.Workspace.SubscriptionId {
			t.Fatalf("Expected %q but got %q for Workspace.SubscriptionId", v.Expected.Workspace.SubscriptionId, actual.Workspace.SubscriptionId)
		}

		if actual.Workspace.ResourceGroup != v.Expected.Workspace.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Workspace.ResourceGroup", v.Expected.Workspace.ResourceGroup, actual.Workspace.ResourceGroup)
		}

		if actual.Workspace.Name != v.Expected.Workspace.Name {
			t.Fatalf("Expected %q but got %q for Workspace.Name", v.Expected.Workspace.Name, actual.Workspace.Name)
		}
	}
}
