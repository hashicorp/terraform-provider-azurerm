package validate

import "testing"

func TestAccountName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{
			name:  "Invalid Character 1",
			input: "DC\\",
			valid: false,
		},
		{
			name:  "Invalid Character 2",
			input: "[abc]",
			valid: false,
		},
		{
			name:  "Valid Account Name",
			input: "acc-test",
			valid: true,
		},
		{
			name:  "Invalid Character 3",
			input: "test&",
			valid: false,
		},
		{
			name:  "Too Few Character",
			input: "ab",
			valid: false,
		},
		{
			name:  "Valid Account Name 2",
			input: "aa-BB_88",
			valid: true,
		},
		{
			name:  "Valid Account Name 3",
			input: "aac-",
			valid: true,
		},
	}
	validationFunction := AccountName()
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
