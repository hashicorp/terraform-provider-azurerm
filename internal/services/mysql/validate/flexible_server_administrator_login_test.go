package validate

import (
	"strings"
	"testing"
)

func TestFlexibleServerAdministratorLogin(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{
			input: "test",
			valid: true,
		},
		{
			input: "administrator",
			valid: false,
		},
		{
			input: "azure_superuser",
			valid: false,
		},
		{
			input: "a",
			valid: true,
		},
		{
			input: "",
			valid: false,
		},
		{
			input: "test_",
			valid: false,
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
			_, err := FlexibleServerAdministratorLogin(tt.input, "administrator_login")
			valid := err == nil
			if valid != tt.valid {
				t.Errorf("Expected valid status %t but got %t for input %s", tt.valid, valid, tt.input)
			}
		})
	}
}
