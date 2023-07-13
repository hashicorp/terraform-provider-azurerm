// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"testing"
)

func TestCostManagementExportID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Error    bool
		Expected *CostManagementExportId
	}{
		{
			Name:  "empty",
			Input: "",
			Error: true,
		},
		{
			Name:  "resource group cost management export",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/providers/Microsoft.CostManagement/exports/export1",
			Expected: &CostManagementExportId{
				Name:  "export1",
				Scope: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo",
			},
		},
		{
			Name:  "resource group cost management export but no name",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/providers/Microsoft.CostManagement/exports/",
			Error: true,
		},
		{
			Name:  "subscription cost management export",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.CostManagement/exports/export1",
			Expected: &CostManagementExportId{
				Name:  "export1",
				Scope: "/subscriptions/00000000-0000-0000-0000-000000000000",
			},
		},
		{
			Name:  "subscription cost management export but no name",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.CostManagement/exports/",
			Error: true,
		},
		{
			Name:  "management group cost management export",
			Input: "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000/providers/Microsoft.CostManagement/exports/export1",
			Expected: &CostManagementExportId{
				Name:  "export1",
				Scope: "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000",
			},
		},
		{
			Name:  "management group cost management export but no name",
			Input: "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000/providers/Microsoft.CostManagement/exports/",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := CostManagementExportID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expected a value but got an error: %+v", err)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q", v.Expected.Name, actual.Name)
		}

		if v.Expected.Scope != actual.Scope {
			t.Fatalf("Expected %+v but got %+v", v.Expected.Scope, actual.Scope)
		}
	}
}
