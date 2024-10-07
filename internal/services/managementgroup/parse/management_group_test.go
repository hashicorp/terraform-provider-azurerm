// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import "testing"

func TestManagementGroupID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Error    bool
		Expected *ManagementGroupId
	}{
		{
			Name:  "Empty",
			Input: "",
			Error: true,
		},
		{
			Name:  "No Management Groups Segment",
			Input: "/providers/Microsoft.Management",
			Error: true,
		},
		{
			Name:  "No Management Group ID",
			Input: "/providers/Microsoft.Management/managementGroups/",
			Error: true,
		},
		{
			Name:  "Management Group ID in UUID",
			Input: "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000",
			Expected: &ManagementGroupId{
				Name: "00000000-0000-0000-0000-000000000000",
			},
		},
		{
			Name:  "Management Group ID in Readable ID",
			Input: "/providers/Microsoft.Management/managementGroups/myGroup",
			Expected: &ManagementGroupId{
				Name: "myGroup",
			},
		},
		{
			Name:  "Management Group ID in Readable ID",
			Input: "/providers/Microsoft.Management/ManagementGroups/myGroup",
			Expected: &ManagementGroupId{
				Name: "myGroup",
			},
		},
		{
			Name:  "Management Group ID in UUID with wrong casing",
			Input: "/providers/microsoft.management/managementgroups/00000000-0000-0000-0000-000000000000",
			Expected: &ManagementGroupId{
				Name: "00000000-0000-0000-0000-000000000000",
			},
		},
		{
			Name:  "Management Group ID in UUID with wrong casing",
			Input: "/providers/microsoft.management/Managementgroups/00000000-0000-0000-0000-000000000000",
			Expected: &ManagementGroupId{
				Name: "00000000-0000-0000-0000-000000000000",
			},
		},
		{
			Name:  "Management Group ID in Readable ID with wrong casing",
			Input: "/providers/microsoft.management/managementgroups/group1",
			Expected: &ManagementGroupId{
				Name: "group1",
			},
		},
		{
			Name:  "Invalid Management group id",
			Input: "/providers/Microsoft.Management/managementGroups/myGroup/another",
			Error: true,
		},
		{
			Name:  "Resource ID in management group",
			Input: "/providers/Microsoft.Management/managementGroups/myGroup/providers/Microsoft.Authorization/policyDefinitions/def1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := ManagementGroupID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}
	}
}

func TestManagementGroupIDForSystemTopic(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Error    bool
		Expected *ManagementGroupId
	}{
		{
			Name:  "Management Group ID for System Topic",
			Input: "/tenants/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000",
			Expected: &ManagementGroupId{
				Name:     "00000000-0000-0000-0000-000000000000",
				TenantID: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := TenantScopedManagementGroupID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}

		if actual.TenantID != v.Expected.TenantID {
			t.Fatalf("Expected %q but got %q for TenantID", v.Expected.TenantID, actual.TenantID)
		}
	}
}
