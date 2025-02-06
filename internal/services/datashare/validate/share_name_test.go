// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestShareName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{
			name:  "invalid character",
			input: "9()",
			valid: false,
		},
		{
			name:  "less character",
			input: "a",
			valid: false,
		},
		{
			name:  "valid",
			input: "adgeFG-98",
			valid: true,
		},
		{
			name:  "valid 2",
			input: "dfakF88u7_",
			valid: true,
		},
	}
	validationFunction := ShareName()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := validationFunction(tt.input, "")
			valid := err == nil
			if valid != tt.valid {
				t.Errorf("expected valid status %t but got %t for input %s", tt.valid, valid, tt.input)
			}
		})
	}
}
