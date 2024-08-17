// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"testing"
)

func TestParseRoleDefinitionID(t *testing.T) {
	testData := []struct {
		RoleDefinitionID string
		Expect           bool
	}{
		{
			RoleDefinitionID: "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleDefinitions/00000000-0000-0000-0000-000000000000|/subscriptions/00000000-0000-0000-0000-000000000000",
			Expect:           true,
		},
		{
			RoleDefinitionID: "/providers/Microsoft.Authorization/roleDefinitions/ab65ee35-dc39-43ac-86a0-accaadf4abcd|/providers/Microsoft.Management/managementGroups/Some-Management-Group",
			Expect:           true,
		},
		{
			RoleDefinitionID: "12345678-1234-1234-1234-1234567890ab|/providers/Microsoft.Management/managementGroups/Some-Management-Group",
			Expect:           false,
		},
	}

	for _, v := range testData {
		t.Logf("parsing %q, expecting %+v", v.RoleDefinitionID, v.Expect)

		roleDefinitionID, err := RoleDefinitionId(v.RoleDefinitionID)
		if err != nil {
			if v.Expect == true {
				t.Fatalf("expected %q parse success, got %q", v.RoleDefinitionID, err)
			}
			t.Logf("parse failed as expected, got %q", err)
		} else {
			if v.Expect == false {
				t.Fatalf("expected %q parse failure, got %q", v.RoleDefinitionID, roleDefinitionID)
			}
			t.Logf("parse succeeded as expected, got %+v", roleDefinitionID)
		}
	}
}
