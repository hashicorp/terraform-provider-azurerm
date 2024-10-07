// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2022-04-01/roleassignments"
)

func TestScopedRoleAssignmentID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *ScopedRoleAssignmentId
	}{
		{
			Input: "",
			Error: true,
		},

		{
			Input: "/",
			Error: true,
		},

		{
			Input: "/providers/Microsoft.Marketplace/providers/Microsoft.Authorization/roleAssignments/",
			Error: true,
		},

		{
			Input: "/providers/Microsoft.Subscription/providers/Microsoft.Authorization/roleAssignments/23456781-2349-8764-5631-234567890121",
			Expected: &ScopedRoleAssignmentId{
				ScopedId: roleassignments.NewScopedRoleAssignmentID("/providers/Microsoft.Subscription", "23456781-2349-8764-5631-234567890121"),
			},
		},

		{
			Input: "/providers/Microsoft.Marketplace/providers/Microsoft.Authorization/roleAssignments/23456781-2349-8764-5631-234567890121",
			Expected: &ScopedRoleAssignmentId{
				ScopedId: roleassignments.NewScopedRoleAssignmentID("/providers/Microsoft.Marketplace", "23456781-2349-8764-5631-234567890121"),
			},
		},

		{
			Input: "/providers/Microsoft.Marketplace/providers/Microsoft.Authorization/roleAssignments/23456781-2349-8764-5631-234567890121|12345678-1234-5678-1234-567890123456",
			Expected: &ScopedRoleAssignmentId{
				ScopedId: roleassignments.NewScopedRoleAssignmentID("/providers/Microsoft.Marketplace", "23456781-2349-8764-5631-234567890121"),
				TenantId: "12345678-1234-5678-1234-567890123456",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ScopedRoleAssignmentID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("expected a value but got an error: %+v", err)
		}

		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

		if actual.ScopedId.RoleAssignmentName != v.Expected.ScopedId.RoleAssignmentName {
			t.Fatalf("Expected %q but got %q for Role Assignment Name", v.Expected.ScopedId.RoleAssignmentName, actual.ScopedId.RoleAssignmentName)
		}

		if actual.TenantId != v.Expected.TenantId {
			t.Fatalf("Expected %q but got %q for Tenant ID", v.Expected.TenantId, actual.TenantId)
		}
	}
}

func TestValidateScopedRoleAssignmentID(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{

		{
			Input: "",
			Valid: false,
		},

		{
			Input: "/",
			Valid: false,
		},

		{
			Input: "/providers/Microsoft.Marketplace/providers/Microsoft.Authorization/roleAssignments/",
			Valid: false,
		},

		{
			Input: "/providers/Microsoft.Subscription/providers/Microsoft.Authorization/roleAssignments/23456781-2349-8764-5631-234567890121",
			Valid: true,
		},

		{
			Input: "/providers/Microsoft.Marketplace/providers/Microsoft.Authorization/roleAssignments/23456781-2349-8764-5631-234567890121",
			Valid: true,
		},

		{
			Input: "/providers/Microsoft.Marketplace/providers/Microsoft.Authorization/roleAssignments/23456781-2349-8764-5631-234567890121|12345678-1234-5678-1234-567890123456",
			Valid: true,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := ValidateScopedRoleAssignmentID(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
