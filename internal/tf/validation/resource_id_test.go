// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validation

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

func TestAsGeneratedID(t *testing.T) {
	cases := []struct {
		input string
		valid bool
	}{
		{
			// canonical
			input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/diskEncryptionSets/set1",
			valid: true,
		},
		{
			// all-lowercase `resourcegroups` was accepted by the legacy parser
			input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/group1/providers/Microsoft.Compute/diskEncryptionSets/set1",
			valid: true,
		},
		{
			// the resource provider namespace was not case-validated by the legacy parser
			input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/microsoft.compute/diskEncryptionSets/set1",
			valid: true,
		},
		{
			input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/group1/providers/MICROSOFT.COMPUTE/diskEncryptionSets/set1",
			valid: true,
		},
		{
			// any other casing of `resourceGroups` was rejected by the legacy parser
			input: "/subscriptions/00000000-0000-0000-0000-000000000000/ResourceGroups/group1/providers/Microsoft.Compute/diskEncryptionSets/set1",
			valid: false,
		},
		{
			// the resource type segment was case-sensitive in the legacy parser
			input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/diskencryptionsets/set1",
			valid: false,
		},
		{
			// `subscriptions` was case-sensitive in the legacy parser
			input: "/Subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/diskEncryptionSets/set1",
			valid: false,
		},
		{
			// wrong resource type entirely
			input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/disks/disk1",
			valid: false,
		},
	}

	validateFunc := AsGeneratedID(commonids.ParseDiskEncryptionSetIDInsensitively)
	for _, tc := range cases {
		_, errors := validateFunc(tc.input, "test")
		if valid := len(errors) == 0; valid != tc.valid {
			t.Errorf("expected valid=%t for %q, got %t (errors: %v)", tc.valid, tc.input, valid, errors)
		}
	}
}
