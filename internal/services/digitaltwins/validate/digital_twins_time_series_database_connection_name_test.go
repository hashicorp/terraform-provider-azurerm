// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"testing"
)

func TestDigitalTwinsTimeSeriesDatabaseConnectionName(t *testing.T) {
	tests := []struct {
		Name  string
		Input string
		Valid bool
	}{
		{
			Name:  "Valid Name",
			Input: "Digital-12-Twins",
			Valid: true,
		},
		{
			Name:  "Min characters",
			Input: "aaa",
			Valid: true,
		},
		{
			Name:  "Max characters",
			Input: "A1234567890123456789012345678901234567890123456789",
			Valid: true,
		},
		{
			Name:  "Empty",
			Input: "",
			Valid: false,
		},
		{
			Name:  "Invalid character",
			Input: "digital_twins",
			Valid: false,
		},
		{
			Name:  "Numbers only",
			Input: "1234",
			Valid: false,
		},
		{
			Name:  "End with `-`",
			Input: "Digital-",
			Valid: false,
		},
		{
			Name:  "Start with `-`",
			Input: "-Digital",
			Valid: false,
		},
		{
			Name:  "Too short",
			Input: "aa",
			Valid: false,
		},
		{
			Name:  "Too long",
			Input: "A12345678901234567890123456789012345678901234567890",
			Valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			_, err := DigitalTwinsTimeSeriesDatabaseConnectionName(tt.Input, "Name")
			valid := err == nil
			if valid != tt.Valid {
				t.Errorf("Expected valid status %t but got %t for input %s", tt.Valid, valid, tt.Input)
			}
		})
	}
}
