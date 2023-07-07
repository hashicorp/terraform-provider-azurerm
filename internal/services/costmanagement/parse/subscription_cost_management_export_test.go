// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = SubscriptionCostManagementExportId{}

func TestSubscriptionCostManagementExportIDFormatter(t *testing.T) {
	actual := NewSubscriptionCostManagementExportID("12345678-1234-9876-4563-123456789012", "export1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.CostManagement/exports/export1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestSubscriptionCostManagementExportID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *SubscriptionCostManagementExportId
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
			// missing ExportName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.CostManagement/",
			Error: true,
		},

		{
			// missing value for ExportName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.CostManagement/exports/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.CostManagement/exports/export1",
			Expected: &SubscriptionCostManagementExportId{
				SubscriptionId: "12345678-1234-9876-4563-123456789012",
				ExportName:     "export1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/PROVIDERS/MICROSOFT.COSTMANAGEMENT/EXPORTS/EXPORT1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := SubscriptionCostManagementExportID(v.Input)
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
		if actual.ExportName != v.Expected.ExportName {
			t.Fatalf("Expected %q but got %q for ExportName", v.Expected.ExportName, actual.ExportName)
		}
	}
}
