// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"testing"
)

func TestPostgresqlFlexibleServerDatabaseCollation(t *testing.T) {
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
			input: "en_US%",
			valid: false,
		},
		{
			name:  "Basic",
			input: "en_US",
			valid: true,
		},
		{
			name:  "Basic",
			input: "en_US.utf8",
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := PostgresqlFlexibleServerDatabaseCollation(tt.input, "collation")
			valid := err == nil
			if valid != tt.valid {
				t.Errorf("Expected valid status %t but got %t for input %s", tt.valid, valid, tt.input)
			}
		})
	}
}
