// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"strings"
	"testing"
)

func TestRoleName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{
			input: "",
			valid: false,
		},
		{
			input: "A",
			valid: false,
		},
		{
			input: "a",
			valid: true,
		},
		{
			input: "a-b",
			valid: false,
		},
		{
			input: "9",
			valid: true,
		},
		{
			input: "a9d",
			valid: true,
		},
		{
			input: strings.Repeat("s", 62),
			valid: true,
		},
		{
			input: strings.Repeat("s", 63),
			valid: true,
		},
		{
			input: strings.Repeat("s", 64),
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := RoleName(tt.input, "name")
			valid := err == nil
			if valid != tt.valid {
				t.Errorf("Expected valid status %t but got %t for input %s", tt.valid, valid, tt.input)
			}
		})
	}
}
