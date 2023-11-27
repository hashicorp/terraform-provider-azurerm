// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = AutomationJobScheduleId{}

func TestAutomationJobScheduleIDFormatter(t *testing.T) {
	actual := NewAutomationJobScheduleID("12345678-1234-9876-4563-123456789012", "group1", "account1", "book1", "schedule1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Automation/automationAccounts/account1/runBook/book1/schedule/schedule1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestAutomationJobScheduleID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *AutomationJobScheduleId
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
			// missing AutomationAccountName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Automation/",
			Error: true,
		},

		{
			// missing value for AutomationAccountName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Automation/automationAccounts/",
			Error: true,
		},

		{
			// missing RunBookName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Automation/automationAccounts/account1/",
			Error: true,
		},

		{
			// missing value for RunBookName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Automation/automationAccounts/account1/runBook/",
			Error: true,
		},

		{
			// missing ScheduleName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Automation/automationAccounts/account1/runBook/book1/",
			Error: true,
		},

		{
			// missing value for ScheduleName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Automation/automationAccounts/account1/runBook/book1/schedule/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Automation/automationAccounts/account1/runBook/book1/schedule/schedule1",
			Expected: &AutomationJobScheduleId{
				SubscriptionId:        "12345678-1234-9876-4563-123456789012",
				ResourceGroup:         "group1",
				AutomationAccountName: "account1",
				RunBookName:           "book1",
				ScheduleName:          "schedule1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/GROUP1/PROVIDERS/MICROSOFT.AUTOMATION/AUTOMATIONACCOUNTS/ACCOUNT1/RUNBOOK/BOOK1/SCHEDULE/SCHEDULE1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := AutomationJobScheduleID(v.Input)
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
		if actual.AutomationAccountName != v.Expected.AutomationAccountName {
			t.Fatalf("Expected %q but got %q for AutomationAccountName", v.Expected.AutomationAccountName, actual.AutomationAccountName)
		}
		if actual.RunBookName != v.Expected.RunBookName {
			t.Fatalf("Expected %q but got %q for RunBookName", v.Expected.RunBookName, actual.RunBookName)
		}
		if actual.ScheduleName != v.Expected.ScheduleName {
			t.Fatalf("Expected %q but got %q for ScheduleName", v.Expected.ScheduleName, actual.ScheduleName)
		}
	}
}
