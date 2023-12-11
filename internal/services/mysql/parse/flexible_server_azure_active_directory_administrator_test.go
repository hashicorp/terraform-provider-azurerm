// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = FlexibleServerAzureActiveDirectoryAdministratorId{}

func TestFlexibleServerAzureActiveDirectoryAdministratorIDFormatter(t *testing.T) {
	actual := NewFlexibleServerAzureActiveDirectoryAdministratorID("12345678-1234-9876-4563-123456789012", "resGroup1", "server1", "ActiveDirectory").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DBforMySQL/flexibleServers/server1/administrators/ActiveDirectory"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestFlexibleServerAzureActiveDirectoryAdministratorID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *FlexibleServerAzureActiveDirectoryAdministratorId
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
			// missing FlexibleServerName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DBforMySQL/",
			Error: true,
		},

		{
			// missing value for FlexibleServerName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DBforMySQL/flexibleServers/",
			Error: true,
		},

		{
			// missing AdministratorName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DBforMySQL/flexibleServers/server1/",
			Error: true,
		},

		{
			// missing value for AdministratorName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DBforMySQL/flexibleServers/server1/administrators/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DBforMySQL/flexibleServers/server1/administrators/ActiveDirectory",
			Expected: &FlexibleServerAzureActiveDirectoryAdministratorId{
				SubscriptionId:     "12345678-1234-9876-4563-123456789012",
				ResourceGroup:      "resGroup1",
				FlexibleServerName: "server1",
				AdministratorName:  "ActiveDirectory",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.DBFORMYSQL/FLEXIBLESERVERS/SERVER1/ADMINISTRATORS/ACTIVEDIRECTORY",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := FlexibleServerAzureActiveDirectoryAdministratorID(v.Input)
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
		if actual.FlexibleServerName != v.Expected.FlexibleServerName {
			t.Fatalf("Expected %q but got %q for FlexibleServerName", v.Expected.FlexibleServerName, actual.FlexibleServerName)
		}
		if actual.AdministratorName != v.Expected.AdministratorName {
			t.Fatalf("Expected %q but got %q for AdministratorName", v.Expected.AdministratorName, actual.AdministratorName)
		}
	}
}
