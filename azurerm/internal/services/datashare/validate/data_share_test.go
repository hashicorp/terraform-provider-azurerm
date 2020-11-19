package validate

import "testing"

func TestDataShareAccountName(t *testing.T) {
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
	validationFunction := DataShareAccountName()
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

func TestDatashareName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{
			name:  "invalid character",
			input: "9()",
			valid: false,
		},
		{
			name:  "less character",
			input: "a",
			valid: false,
		},
		{
			name:  "invalid character2",
			input: "adgeFG-98",
			valid: false,
		},
		{
			name:  "valid",
			input: "dfakF88u7_",
			valid: true,
		},
	}
	validationFunction := DatashareName()
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

func TestDatashareSyncName(t *testing.T) {
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
	validationFunction := DataShareSyncName()
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
