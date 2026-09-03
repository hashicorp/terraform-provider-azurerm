// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"testing"
)

func TestValidateCognitiveServicesRaiBlocklistName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{
			name:  "empty name",
			input: "",
			valid: false,
		},
		{
			name:  "Invalid single character",
			input: "a",
			valid: false,
		},
		{
			name:  "Valid two characters",
			input: "ab",
			valid: true,
		},
		{
			name:  "Valid alphanumeric",
			input: "hello1",
			valid: true,
		},
		{
			name:  "Valid contains underscore",
			input: "hello_world",
			valid: true,
		},
		{
			name:  "Valid contains hyphen",
			input: "hello-world",
			valid: true,
		},
		{
			name:  "Valid starts with digit",
			input: "1hello",
			valid: true,
		},
		{
			name:  "Valid starts with underscore",
			input: "_hello",
			valid: true,
		},
		{
			name:  "Valid starts with hyphen",
			input: "-hello",
			valid: true,
		},
		{
			name:  "Valid ends with underscore",
			input: "hello_",
			valid: true,
		},
		{
			name:  "Valid ends with hyphen",
			input: "hello-",
			valid: true,
		},
		{
			name:  "Invalid contains period",
			input: "a.bc",
			valid: false,
		},
		{
			name:  "Invalid contains whitespace",
			input: "a bc",
			valid: false,
		},
		{
			name:  "Invalid contains special character",
			input: "a@bc",
			valid: false,
		},
		{
			name:  "contains slash",
			input: "a/bc",
			valid: false,
		},
		{
			name:  "Valid 64 characters",
			input: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789012",
			valid: true,
		},
		{
			name:  "Invalid starts with period",
			input: ".heyo",
			valid: false,
		},
		{
			name:  "Invalid 65 characters",
			input: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890123",
			valid: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := RaiBlocklistName()(tt.input, "")
			valid := err == nil
			if valid != tt.valid {
				t.Errorf("Expected valid status %t but got %t for input %s", tt.valid, valid, tt.input)
			}
		})
	}
}
