// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"testing"
)

func TestIntegrationAccountAssemblyVersion(t *testing.T) {
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
			input: "1.1",
			valid: true,
		},
		{
			input: "1.1.1",
			valid: false,
		},
		{
			input: "1.1.1.1",
			valid: true,
		},
		{
			input: "1.1.1.1.1",
			valid: false,
		},
	}

	validationFunction := IntegrationAccountAssemblyVersion()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := validationFunction(tt.input, "")
			valid := err == nil
			if valid != tt.valid {
				t.Errorf("Expected valid status %t but got %t for input %s", tt.valid, valid, tt.input)
			}
		})
	}
}
