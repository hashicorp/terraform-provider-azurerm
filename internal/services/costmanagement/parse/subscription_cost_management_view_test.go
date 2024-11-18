// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = SubscriptionCostManagementViewId{}

func TestSubscriptionCostManagementViewIDFormatter(t *testing.T) {
	actual := NewSubscriptionCostManagementViewID("12345678-1234-9876-4563-123456789012", "view1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.CostManagement/views/view1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestSubscriptionCostManagementViewID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *SubscriptionCostManagementViewId
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
			// missing ViewName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.CostManagement/",
			Error: true,
		},

		{
			// missing value for ViewName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.CostManagement/views/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.CostManagement/views/view1",
			Expected: &SubscriptionCostManagementViewId{
				SubscriptionId: "12345678-1234-9876-4563-123456789012",
				ViewName:       "view1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/PROVIDERS/MICROSOFT.COSTMANAGEMENT/VIEWS/VIEW1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := SubscriptionCostManagementViewID(v.Input)
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
		if actual.ViewName != v.Expected.ViewName {
			t.Fatalf("Expected %q but got %q for ViewName", v.Expected.ViewName, actual.ViewName)
		}
	}
}
