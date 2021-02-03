package eventhub

import (
	"strings"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventhub/validate"
)

func TestValidateEventHubName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{
			name:  "Valid short name",
			input: "abc",
			valid: true,
		},
		{
			name:  "Valid short name",
			input: "a",
			valid: true,
		},
		{
			name:  "Valid name with dot",
			input: "a.b",
			valid: true,
		},
		{
			name:  "Just a digit",
			input: "1",
			valid: true,
		},
		{
			name:  "Invalid long name",
			input: strings.Repeat("a", 51),
			valid: false,
		},
		{
			name:  "Invalid short name",
			input: ".",
			valid: false,
		},
		{
			name:  "Invalid name with period at end",
			input: "a.",
			valid: false,
		},
		{
			name:  "empty name",
			input: "",
			valid: false,
		},
	}
	validationFunction := validate.ValidateEventHubName()
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
