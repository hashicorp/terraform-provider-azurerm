// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package monitor

import "testing"

func TestNormalizeMonitorSeverity(t *testing.T) {
	cases := []struct {
		input       string
		expected    int
		expectError bool
	}{
		{input: "0", expected: 0},
		{input: "1", expected: 1},
		{input: "2", expected: 2},
		{input: "3", expected: 3},
		{input: "4", expected: 4},
		{input: "critical", expected: 0},
		{input: "error", expected: 1},
		{input: "warning", expected: 2},
		{input: "informational", expected: 3},
		{input: "verbose", expected: 4},
		{input: " Warning ", expected: 2},
		{input: "invalid", expectError: true},
	}

	for _, tc := range cases {
		actual, err := normalizeMonitorSeverity(tc.input)
		if tc.expectError {
			if err == nil {
				t.Fatalf("expected error for %q", tc.input)
			}
			continue
		}

		if err != nil {
			t.Fatalf("unexpected error for %q: %s", tc.input, err)
		}

		if actual != tc.expected {
			t.Fatalf("expected %d, got %d for %q", tc.expected, actual, tc.input)
		}
	}
}

func TestNormalizeMonitorSeverityState(t *testing.T) {
	if actual := normalizeMonitorSeverityState("warning"); actual != "2" {
		t.Fatalf("expected warning to normalize to 2, got %q", actual)
	}

	if actual := normalizeMonitorSeverityState(" 3 "); actual != "3" {
		t.Fatalf("expected 3 to normalize to 3, got %q", actual)
	}
}

func TestSuppressMonitorSeverityDiff(t *testing.T) {
	if !suppressMonitorSeverityDiff("severity", "2", "warning", nil) {
		t.Fatal("expected 2 and warning to suppress diff")
	}

	if suppressMonitorSeverityDiff("severity", "1", "warning", nil) {
		t.Fatal("did not expect 1 and warning to suppress diff")
	}
}
