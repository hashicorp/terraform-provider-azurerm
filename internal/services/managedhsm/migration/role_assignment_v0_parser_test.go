package migration

import (
	"testing"
)

func TestRoleAssignmentV0Parser(t *testing.T) {
	testData := []struct {
		input    string
		expected *legacyV0RoleAssignmentId
	}{
		{
			input:    "",
			expected: nil,
		},
		{
			// missing scope
			input:    "https://my-hsm.managedhsm.azure.net/RoleAssignment/test",
			expected: nil,
		},
		{
			// missing role assignment name
			input:    "https://my-hsm.managedhsm.azure.net///RoleAssignment/",
			expected: nil,
		},
		{
			// wrong legacy type
			input:    "https://my-hsm.managedhsm.azure.net///RoleDefinition/example",
			expected: nil,
		},
		{
			// Public
			input: "https://my-hsm.managedhsm.azure.net///RoleAssignment/test",
			expected: &legacyV0RoleAssignmentId{
				managedHSMName:     "my-hsm",
				domainSuffix:       "managedhsm.azure.net",
				scope:              "/",
				roleAssignmentName: "test",
			},
		},
		{
			// Public - superfluous port
			input: "https://my-hsm.managedhsm.azure.net:443///RoleAssignment/test",
			expected: &legacyV0RoleAssignmentId{
				managedHSMName:     "my-hsm",
				domainSuffix:       "managedhsm.azure.net",
				scope:              "/",
				roleAssignmentName: "test",
			},
		},
		{
			input: "https://my-hsm.managedhsm.azure.cn///RoleAssignment/test",
			expected: &legacyV0RoleAssignmentId{
				managedHSMName:     "my-hsm",
				domainSuffix:       "managedhsm.azure.cn",
				scope:              "/",
				roleAssignmentName: "test",
			},
		},
		{
			input: "https://my-hsm.managedhsm.usgovcloudapi.net///RoleAssignment/test",
			expected: &legacyV0RoleAssignmentId{
				managedHSMName:     "my-hsm",
				domainSuffix:       "managedhsm.usgovcloudapi.net",
				scope:              "/",
				roleAssignmentName: "test",
			},
		},
	}
	for _, test := range testData {
		t.Logf("Testing %q..", test.input)
		actual, err := parseLegacyV0RoleAssignmentId(test.input)
		if err != nil {
			if test.expected == nil {
				continue
			}

			t.Fatalf("unexpected error: %+v", err)
		}

		if test.expected == nil {
			if actual == nil {
				continue
			}

			t.Fatalf("expected nothing but got %+v", actual)
		}
		if actual == nil {
			t.Fatalf("expected %+v but got nil", test.expected)
		}
	}
}
