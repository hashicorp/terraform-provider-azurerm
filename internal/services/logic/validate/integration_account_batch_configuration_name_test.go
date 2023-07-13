// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"strings"
	"testing"
)

func TestIntegrationAccountBatchConfigurationName(t *testing.T) {
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
			input: "test1",
			valid: true,
		},
		{
			input: "a2-.()b",
			valid: false,
		},
		{
			input: "a2&b",
			valid: false,
		},
		{
			input: strings.Repeat("s", 19),
			valid: true,
		},
		{
			input: strings.Repeat("s", 20),
			valid: true,
		},
		{
			input: strings.Repeat("s", 21),
			valid: false,
		},
	}

	validationFunction := IntegrationAccountBatchConfigurationName()
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
