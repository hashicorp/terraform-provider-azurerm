package validate

import "testing"

func TestSnapshotScheduleName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{
			name:  "valid characters",
			input: "*() _-$@!",
			valid: true,
		},
		{
			name:  "Empty",
			input: "",
			valid: false,
		},
		{
			name:  "invalid characters",
			input: "&^*",
			valid: false,
		},
		{
			name:  "invalid characters",
			input: "dfwe%",
			valid: false,
		},
	}
	validationFunction := SnapshotScheduleName()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := validationFunction(tt.input, "")
			valid := err == nil
			if valid != tt.valid {
				t.Errorf("expected valid status %t but got %t for input %s", tt.valid, valid, tt.input)
			}
		})
	}
}
