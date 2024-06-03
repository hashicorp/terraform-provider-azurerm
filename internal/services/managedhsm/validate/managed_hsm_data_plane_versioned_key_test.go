// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestManagedHSMDataPlaneVersionedKeyID(t *testing.T) {
	testData := []struct {
		input string
		valid bool
	}{
		{
			// empty = invalid
			input: "",
			valid: false,
		},
		{
			// key vault versioned key id
			input: "https://example.keyvault.azure.net/keys/abc123/bcd234",
			valid: false,
		},
		{
			// domain but no uri
			input: "https://example.keyvault.azure.net/",
			valid: false,
		},
		{
			// managed hsm domain but wrong type
			input: "https://example.managedhsm.azure.net/numbers/abc123/bcd234",
			valid: false,
		},
		{
			// managed hsm key id (no version)
			input: "https://example.managedhsm.azure.net/keys/abc123",
			valid: false,
		},
		{
			// managed hsm key id (with version)
			input: "https://example.managedhsm.azure.net/keys/abc123/bcd234",
			valid: true,
		},
		{
			// managed hsm key id (with version but extra)
			input: "https://example.managedhsm.azure.net/keys/abc123/bcd234/cde345",
			valid: false,
		},
	}
	for _, item := range testData {
		t.Logf("Testing %q", item.input)
		warnings, errs := ManagedHSMDataPlaneVersionedKeyID(item.input, "some_field")
		actual := len(warnings) == 0 && len(errs) == 0
		if item.valid != actual {
			t.Fatalf("expected %t but got %t", item.valid, actual)
		}
	}
}
