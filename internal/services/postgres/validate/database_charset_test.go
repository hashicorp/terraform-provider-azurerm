// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestDatabaseCharset(t *testing.T) {
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
			name:  "Invalid",
			input: "UTF-8",
			valid: false,
		},
		{
			name:  "Uppercase",
			input: "UTF8",
			valid: true,
		},
		{
			name:  "Lowercase",
			input: "utf8",
			valid: true,
		},
		{
			name:  "With underscore",
			input: "SQL_ASCII",
			valid: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DatabaseCharset(tt.input, "charset")
			valid := err == nil
			if valid != tt.valid {
				t.Errorf("Expected valid status %t but got %t for input %s", tt.valid, valid, tt.input)
			}
		})
	}
}
