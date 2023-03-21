package validate

import (
	"strings"
	"testing"
)

func TestAdministratorLoginPasswordName(t *testing.T) {
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
			input: strings.Repeat("s", 255),
			valid: true,
		},
		{
			input: strings.Repeat("s", 256),
			valid: true,
		},
		{
			input: strings.Repeat("s", 257),
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := AdministratorLoginPassword(tt.input, "administrator_login_password")
			valid := err == nil
			if valid != tt.valid {
				t.Errorf("Expected valid status %t but got %t for input %s", tt.valid, valid, tt.input)
			}
		})
	}
}
