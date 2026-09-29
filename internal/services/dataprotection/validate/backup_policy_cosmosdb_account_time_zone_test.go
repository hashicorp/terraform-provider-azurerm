// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestBackupPolicyCosmosdbAccountTimeZone(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{
			input:    "UTC",
			expected: true,
		},
		{
			input:    "Coordinated Universal Time",
			expected: true,
		},
		{
			input:    "Invalid",
			expected: false,
		},
	}

	validator := BackupPolicyCosmosdbAccountTimeZone()
	for _, testCase := range testCases {
		_, errors := validator(testCase.input, "time_zone")
		result := len(errors) == 0
		if result != testCase.expected {
			t.Fatalf("expected validation result for %q to be %t, got %t", testCase.input, testCase.expected, result)
		}
	}
}
