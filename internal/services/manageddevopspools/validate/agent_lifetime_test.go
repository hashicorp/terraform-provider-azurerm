// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestAgentLifetime(t *testing.T) {
	testCases := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		{
			name:      "valid hours only",
			input:     "08:00:00",
			shouldErr: false,
		},
		{
			name:      "valid with minutes and seconds",
			input:     "01:30:45",
			shouldErr: false,
		},
		{
			name:      "valid zero",
			input:     "00:00:00",
			shouldErr: false,
		},
		{
			name:      "valid max without days",
			input:     "23:59:59",
			shouldErr: false,
		},
		{
			name:      "valid with days",
			input:     "1.08:00:00",
			shouldErr: false,
		},
		{
			name:      "valid 7 days exactly",
			input:     "7.00:00:00",
			shouldErr: false,
		},
		{
			name:      "valid multi-digit days",
			input:     "3.12:30:00",
			shouldErr: false,
		},
		{
			name:      "empty string",
			input:     "",
			shouldErr: true,
		},
		{
			name:      "ISO 8601 duration",
			input:     "P7D",
			shouldErr: true,
		},
		{
			name:      "missing seconds",
			input:     "08:00",
			shouldErr: true,
		},
		{
			name:      "plain number",
			input:     "3600",
			shouldErr: true,
		},
		{
			name:      "Go duration format",
			input:     "8h0m0s",
			shouldErr: true,
		},
		{
			name:      "hours out of range",
			input:     "24:00:00",
			shouldErr: true,
		},
		{
			name:      "minutes out of range",
			input:     "00:60:00",
			shouldErr: true,
		},
		{
			name:      "seconds out of range",
			input:     "00:00:60",
			shouldErr: true,
		},
		{
			name:      "exceeds 7 days",
			input:     "8.00:00:00",
			shouldErr: true,
		},
		{
			name:      "exactly one second over 7 days",
			input:     "7.00:00:01",
			shouldErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, errors := AgentLifetime(tc.input, "test")

			if tc.shouldErr && len(errors) == 0 {
				t.Errorf("Expected validation to fail for input %q, but it passed", tc.input)
			}

			if !tc.shouldErr && len(errors) > 0 {
				t.Errorf("Expected validation to pass for input %q, but it failed with errors: %v", tc.input, errors)
			}
		})
	}
}
