// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"strings"
	"testing"
)

func TestCapabilityHostName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{
			name:  "empty",
			input: "",
			valid: false,
		},
		{
			name:  "single character",
			input: "a",
			valid: true,
		},
		{
			name:  "maximum length",
			input: "a" + strings.Repeat("b", 254),
			valid: true,
		},
		{
			name:  "too long",
			input: "a" + strings.Repeat("b", 255),
			valid: false,
		},
		{
			name:  "starts with digit",
			input: "1example",
			valid: true,
		},
		{
			name:  "starts with dash",
			input: "-example",
			valid: false,
		},
		{
			name:  "starts with underscore",
			input: "_example",
			valid: false,
		},
		{
			name:  "contains dash and underscore",
			input: "example-capability_host",
			valid: true,
		},
		{
			name:  "contains period",
			input: "example.capability-host",
			valid: false,
		},
		{
			name:  "contains space",
			input: "example capability host",
			valid: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, errors := CapabilityHostName()(test.input, "name")
			if (len(errors) == 0) != test.valid {
				t.Errorf("expected valid status %t but got %t for input %q", test.valid, len(errors) == 0, test.input)
			}
		})
	}
}
