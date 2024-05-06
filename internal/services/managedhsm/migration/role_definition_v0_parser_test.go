// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"testing"
)

func TestRoleDefinitionV0Parser(t *testing.T) {
	testData := []struct {
		input    string
		expected *legacyV0RoleDefinitionId
	}{
		{
			input:    "",
			expected: nil,
		},
		{
			// missing scope
			input:    "https://my-hsm.managedhsm.azure.net/RoleDefinition/test",
			expected: nil,
		},
		{
			// missing role assignment name
			input:    "https://my-hsm.managedhsm.azure.net///RoleDefinition/",
			expected: nil,
		},
		{
			// wrong legacy type
			input:    "https://my-hsm.managedhsm.azure.net///RoleAssignment/example",
			expected: nil,
		},
		{
			// Public
			input: "https://my-hsm.managedhsm.azure.net///RoleDefinition/test",
			expected: &legacyV0RoleDefinitionId{
				managedHSMName:     "my-hsm",
				domainSuffix:       "managedhsm.azure.net",
				scope:              "/",
				roleDefinitionName: "test",
			},
		},
		{
			// Public - superfluous port
			input: "https://my-hsm.managedhsm.azure.net:443///RoleDefinition/test",
			expected: &legacyV0RoleDefinitionId{
				managedHSMName:     "my-hsm",
				domainSuffix:       "managedhsm.azure.net",
				scope:              "/",
				roleDefinitionName: "test",
			},
		},
		{
			// Public - invalid port
			input:    "https://my-hsm.managedhsm.azure.net:445///RoleDefinition/test",
			expected: nil,
		},
		{
			input: "https://my-hsm.managedhsm.azure.cn///RoleDefinition/test",
			expected: &legacyV0RoleDefinitionId{
				managedHSMName:     "my-hsm",
				domainSuffix:       "managedhsm.azure.cn",
				scope:              "/",
				roleDefinitionName: "test",
			},
		},
		{
			input: "https://my-hsm.managedhsm.usgovcloudapi.net///RoleDefinition/test",
			expected: &legacyV0RoleDefinitionId{
				managedHSMName:     "my-hsm",
				domainSuffix:       "managedhsm.usgovcloudapi.net",
				scope:              "/",
				roleDefinitionName: "test",
			},
		},
	}
	for _, test := range testData {
		t.Logf("Testing %q..", test.input)
		actual, err := parseLegacyV0RoleDefinitionId(test.input)
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

		if test.expected.managedHSMName != actual.managedHSMName {
			t.Fatalf("expected managedHSMName to be %q but got %q", test.expected.managedHSMName, actual.managedHSMName)
		}
		if test.expected.domainSuffix != actual.domainSuffix {
			t.Fatalf("expected domainSuffix to be %q but got %q", test.expected.domainSuffix, actual.domainSuffix)
		}
		if test.expected.scope != actual.scope {
			t.Fatalf("expected scope to be %q but got %q", test.expected.scope, actual.scope)
		}
		if test.expected.roleDefinitionName != actual.roleDefinitionName {
			t.Fatalf("expected roleDefinitionName to be %q but got %q", test.expected.roleDefinitionName, actual.roleDefinitionName)
		}
	}
}
