// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"strings"
	"testing"
)

func TestFlexibleServerDatabaseName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{
			name:  "Empty",
			input: "",
			valid: false,
		},
		{
			name:  "Invalid Characters",
			input: "flexdb%",
			valid: false,
		},
		{
			name:  "Start with integer",
			input: "1flexdb",
			valid: false,
		},
		{
			name:  "Name to long",
			input: strings.Repeat("a", 64),
			valid: false,
		},
		{
			name:  "One character",
			input: "a",
			valid: true,
		},
		{
			name:  "End with `_`",
			input: "flexdb_",
			valid: true,
		},
		{
			name:  "Start with `_`",
			input: "_flexdb",
			valid: true,
		},
		{
			name:  "Valid",
			input: "flexdb-1-test",
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := FlexibleServerDatabaseName(tt.input, "name")
			valid := err == nil
			if valid != tt.valid {
				t.Errorf("Expected valid status %t but got %t for input %s", tt.valid, valid, tt.input)
			}
		})
	}
}
