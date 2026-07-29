// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package authorization

import "testing"

func TestRoleAssignmentPrincipalType(t *testing.T) {
	schema := resourceArmRoleAssignment().Schema["principal_type"]

	for _, value := range []string{
		"AgentServicePrincipal",
		"AgentUser",
		"Group",
		"ServicePrincipal",
		"User",
	} {
		t.Run(value, func(t *testing.T) {
			_, errors := schema.ValidateFunc(value, "principal_type")
			if len(errors) != 0 {
				t.Fatalf("expected %q to be a valid principal_type, got: %v", value, errors)
			}
		})
	}

	_, errors := schema.ValidateFunc("Unknown", "principal_type")
	if len(errors) == 0 {
		t.Fatal("expected an unknown principal_type to be rejected")
	}
}
