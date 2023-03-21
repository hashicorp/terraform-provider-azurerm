package validate

import (
	"strings"
	"testing"
)

func TestClusterName(t *testing.T) {
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
			input: strings.Repeat("s", 259),
			valid: true,
		},
		{
			input: strings.Repeat("s", 260),
			valid: true,
		},
		{
			input: strings.Repeat("s", 261),
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ClusterName(tt.input, "name")
			valid := err == nil
			if valid != tt.valid {
				t.Errorf("Expected valid status %t but got %t for input %s", tt.valid, valid, tt.input)
			}
		})
	}
}
