package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = SecurityCenterWorkspaceId{}

func TestSecurityCenterWorkspaceIDFormatter(t *testing.T) {
	actual := NewSecurityCenterWorkspaceID("12345678-1234-9876-4563-123456789012", "workspaceSettings1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Security/workspaceSettings/workspaceSettings1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestSecurityCenterWorkspaceID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *SecurityCenterWorkspaceId
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
			// missing WorkspaceSettingName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Security/",
			Error: true,
		},

		{
			// missing value for WorkspaceSettingName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Security/workspaceSettings/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Security/workspaceSettings/workspaceSettings1",
			Expected: &SecurityCenterWorkspaceId{
				SubscriptionId:       "12345678-1234-9876-4563-123456789012",
				WorkspaceSettingName: "workspaceSettings1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/PROVIDERS/MICROSOFT.SECURITY/WORKSPACESETTINGS/WORKSPACESETTINGS1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := SecurityCenterWorkspaceID(v.Input)
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
		if actual.WorkspaceSettingName != v.Expected.WorkspaceSettingName {
			t.Fatalf("Expected %q but got %q for WorkspaceSettingName", v.Expected.WorkspaceSettingName, actual.WorkspaceSettingName)
		}
	}
}
