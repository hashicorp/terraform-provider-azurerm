// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"strings"
	"testing"
)

func TestExascaleDatabaseStorageVaultName(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			input:    "",
			expected: false,
		},
		{
			input:    "v",
			expected: true,
		},
		{
			input:    "goodName",
			expected: true,
		},
		{
			input:    "_",
			expected: true,
		},
		{
			input:    "_good-Name-",
			expected: true,
		},
		{
			input:    "_G0od_name_",
			expected: true,
		},
		{
			input:    "good3Name",
			expected: true,
		},
		{
			input:    "good-name",
			expected: true,
		},
		{
			input:    "_G0od_name_",
			expected: true,
		},
		{
			input:    "Bad-Name2--",
			expected: false,
		},
		{
			input:    "1",
			expected: false,
		},
		{
			input:    "-Bad-Name2",
			expected: false,
		},
		{
			input:    "-Bad-Name2--",
			expected: false,
		},
		{
			input:    "--bad-Name",
			expected: false,
		},
		{
			input:    "2Bad-Name",
			expected: false,
		},
		{
			input:    "another--bad-name",
			expected: false,
		},
		{
			input:    "b@d2name",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := ExascaleDatabaseResourceName(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}

func TestExascaleDatabaseVirtualMachineClusterSSHPublicKeys(t *testing.T) {
	testData := []struct {
		name     string
		input    interface{}
		expected bool
	}{
		{
			name:     "multiple keys",
			input:    []interface{}{"ssh-rsa aaaaaa", "ssh-ed25519 bbbbbbb"},
			expected: true,
		},
		{
			name:     "combined length equal to limit",
			input:    []interface{}{strings.Repeat("a", 5000), strings.Repeat("b", 2500), strings.Repeat("b", 2500)},
			expected: true,
		},
		{
			name:     "combined length exceeding limit",
			input:    []interface{}{strings.Repeat("a", 5000), strings.Repeat("b", 2500), strings.Repeat("b", 2501)},
			expected: false,
		},
		{
			name:     "single key exceeding limit",
			input:    []interface{}{strings.Repeat("a", 10001)},
			expected: false,
		},
		{
			name:     "model string slice",
			input:    []string{"ssh-rsa aaaaaa", "ssh-ed25519 bbbbbb"},
			expected: true,
		},
		{
			name:     "not a list",
			input:    "ssh-rsa AAAA",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Run(v.name, func(t *testing.T) {
			_, errors := ExascaleDatabaseVirtualMachineClusterSSHPublicKeys(v.input, "ssh_public_keys")
			actual := len(errors) == 0
			if v.expected != actual {
				t.Fatalf("expected %t but got %t", v.expected, actual)
			}
		})
	}
}
