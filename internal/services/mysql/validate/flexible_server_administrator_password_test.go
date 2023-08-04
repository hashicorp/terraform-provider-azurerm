// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"strings"
	"testing"
)

func TestFlexibleServerAdministratorPassword(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{
			input: "testtest",
			valid: false,
		},
		{
			input: "aA7",
			valid: false,
		},
		{
			input: "aaaaAAAA7",
			valid: true,
		},
		{
			input: "Aaaaaaaa#",
			valid: true,
		},
		{
			input: "Aaa$aaaaa",
			valid: true,
		},
		{
			input: "%Aaaaaaaa",
			valid: true,
		},
		{
			input: "a7888888@",
			valid: true,
		},
		{
			input: "AAAAAAA7!",
			valid: true,
		},
		{
			input: "abbbbbbb#",
			valid: false,
		},
		{
			input: "A#" + strings.Repeat("s", 125),
			valid: true,
		},
		{
			input: "A#" + strings.Repeat("s", 126),
			valid: true,
		},
		{
			input: "A#" + strings.Repeat("s", 127),
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := FlexibleServerAdministratorPassword(tt.input, "administrator_password")
			valid := err == nil
			if valid != tt.valid {
				t.Errorf("Expected valid status %t but got %t for input %s", tt.valid, valid, tt.input)
			}
		})
	}
}
