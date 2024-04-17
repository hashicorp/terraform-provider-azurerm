// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"testing"
)

func TestManagedHSMDataPlaneRoleAssignmentID(t *testing.T) {
	cases := []struct {
		Input       string
		ExpectError bool
	}{
		{
			Input:       "",
			ExpectError: true,
		},
		{
			Input:       "https://my-hsm.managedhsm.azure.net///test",
			ExpectError: true,
		},
		{
			Input:       "https://my-hsm.managedhsm.azure.net////providers/Microsoft.Authorization/roleAssignments/test",
			ExpectError: false,
		},
		{
			Input:       "https://my-hsm.managedhsm.azure.net//keys//providers/Microsoft.Authorization/roleAssignments/1492",
			ExpectError: false,
		},
		{
			Input:       "https://my-hsm.managedhsm.azure.cn////providers/Microsoft.Authorization/roleAssignments/test",
			ExpectError: false,
		},
		{
			Input:       "https://my-hsm.managedhsm.usgovcloudapi.net//keys//providers/Microsoft.Authorization/roleAssignments/1492",
			ExpectError: false,
		},
		{
			Input:       "https://my-hsm.managedhsm.azure.net//keys//providers/Microsoft.Authorization/roleAssignments/1492/suffix",
			ExpectError: true,
		},
		{
			Input:       "https://my-hsm.managedhsm.azure.net////providers/Microsoft.Authorization/roleDefinitions/000-000",
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		t.Logf("Testing %q..", tc.Input)
		warnings, err := ManagedHSMDataPlaneRoleAssignmentID(tc.Input, "example")
		if err != nil {
			if !tc.ExpectError {
				t.Fatalf("Got error for input %q: %+v", tc.Input, err)
			}

			continue
		}

		if tc.ExpectError && len(warnings) == 0 {
			t.Fatalf("Got no errors for input %q but expected some", tc.Input)
		} else if !tc.ExpectError && len(warnings) > 0 {
			t.Fatalf("Got %d errors for input %q when didn't expect any", len(warnings), tc.Input)
		}
	}
}
