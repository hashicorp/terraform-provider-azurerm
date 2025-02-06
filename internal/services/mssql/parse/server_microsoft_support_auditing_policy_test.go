// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = ServerMicrosoftSupportAuditingPolicyId{}

func TestServerMicrosoftSupportAuditingPolicyIDFormatter(t *testing.T) {
	actual := NewServerMicrosoftSupportAuditingPolicyID("12345678-1234-9876-4563-123456789012", "group1", "server1", "default").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Sql/servers/server1/devOpsAuditingSettings/default"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestServerMicrosoftSupportAuditingPolicyID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *ServerMicrosoftSupportAuditingPolicyId
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
			// missing ServerName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Sql/",
			Error: true,
		},

		{
			// missing value for ServerName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Sql/servers/",
			Error: true,
		},

		{
			// missing DevOpsAuditingSettingName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Sql/servers/server1/",
			Error: true,
		},

		{
			// missing value for DevOpsAuditingSettingName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Sql/servers/server1/devOpsAuditingSettings/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Sql/servers/server1/devOpsAuditingSettings/default",
			Expected: &ServerMicrosoftSupportAuditingPolicyId{
				SubscriptionId:            "12345678-1234-9876-4563-123456789012",
				ResourceGroup:             "group1",
				ServerName:                "server1",
				DevOpsAuditingSettingName: "default",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/GROUP1/PROVIDERS/MICROSOFT.SQL/SERVERS/SERVER1/DEVOPSAUDITINGSETTINGS/DEFAULT",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ServerMicrosoftSupportAuditingPolicyID(v.Input)
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
		if actual.ServerName != v.Expected.ServerName {
			t.Fatalf("Expected %q but got %q for ServerName", v.Expected.ServerName, actual.ServerName)
		}
		if actual.DevOpsAuditingSettingName != v.Expected.DevOpsAuditingSettingName {
			t.Fatalf("Expected %q but got %q for DevOpsAuditingSettingName", v.Expected.DevOpsAuditingSettingName, actual.DevOpsAuditingSettingName)
		}
	}
}
