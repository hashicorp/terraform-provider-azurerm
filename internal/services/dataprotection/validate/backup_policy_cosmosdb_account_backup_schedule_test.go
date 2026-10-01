// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestBackupPolicyCosmosdbAccountBackupSchedule(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{
			input:    "R/2026-02-08T10:00+00:00/P1W",
			expected: true,
		},
		{
			input:    "R/2026-02-08T10:00:00+00:00/P1W",
			expected: true,
		},
		{
			input:    "R/2026-02-08T10:00:00.123Z/P1W",
			expected: true,
		},
		{
			input:    "R/2026-02-08T10:00:00-07:00/P1W",
			expected: true,
		},
		{
			input:    "R5/2026-02-08T10:00:00Z/P1W",
			expected: false,
		},
		{
			input:    "R/20260208T10:00:00Z/P1W",
			expected: false,
		},
		{
			input:    "R/2026-02-08T10Z/P1W",
			expected: false,
		},
		{
			input:    "R/2026-02-08T10:00:00.12Z/P1W",
			expected: false,
		},
		{
			input:    "R/2026-02-08T10:00:00.1234Z/P1W",
			expected: false,
		},
		{
			input:    "R/2026-02-08T10:00:00+0000/P1W",
			expected: false,
		},
		{
			input:    "R/2026-02-08T10:00:00/P1W",
			expected: false,
		},
		{
			input:    "R/2026-02-08T10:00:00Z/P1D",
			expected: false,
		},
		{
			input:    "R/2026-02-08T10:00:00Z/P1W/extra",
			expected: false,
		},
		{
			input:    "R/2023-10-15T14:30:45.123+05:30//P1W",
			expected: false,
		},
	}

	validator := BackupPolicyCosmosdbAccountBackupSchedule()
	for _, testCase := range testCases {
		_, errors := validator(testCase.input, "backup_schedule")
		result := len(errors) == 0
		if result != testCase.expected {
			t.Fatalf("expected validation result for %q to be %t, got %t", testCase.input, testCase.expected, result)
		}
	}
}
