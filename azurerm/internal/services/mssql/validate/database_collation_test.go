package validate

import (
	"testing"
)

func TestDatabaseCollation(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{
			name:  "SQL Collation",
			input: "SQL_Latin1_General_CP1_CI_AS",
			valid: true,
		},
		{
			name:  "Windows Collation",
			input: "Latin1_General_100_CI_AS_SC",
			valid: true,
		},
		{
			name:  "SQL Collation",
			input: "SQL_AltDiction_CP850_CI_AI",
			valid: true,
		},
		{
			name:  "SQL Collation",
			input: "SQL_Croatian_CP1250_CI_AS",
			valid: true,
		},
		{
			name:  "Windows Collation",
			input: "Chinese_Hong_Kong_Stroke_90_CI_AI",
			valid: true,
		},
		{
			name:  "Windows Collation",
			input: "Japanese_BIN",
			valid: true,
		},
		{
			name:  "lowercase",
			input: "sql_croatian_cp1250_ci_as",
			valid: false,
		},
		{
			name:  "extra dot",
			input: "SQL_Croatian_CP1250.",
			valid: false,
		},
		{
			name:  "Invalid collation",
			input: "CDD",
			valid: false,
		},
		{
			name:  "Double definition",
			input: "Latin1_General_100_CI_CS",
			valid: false,
		},
	}
	validationFunction := MsSqlDBCollation()
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
