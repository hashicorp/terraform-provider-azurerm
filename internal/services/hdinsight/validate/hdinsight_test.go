package validate

import "testing"

func TestHDInsightClusterVersion(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{
			name:  "empty name",
			input: "",
			valid: false,
		},
		{
			name:  "major only",
			input: "1",
			valid: false,
		},
		{
			name:  "major minor",
			input: "1.2",
			valid: true,
		},
		{
			name:  "major minor large",
			input: "1000.2000",
			valid: true,
		},
		{
			name:  "major minor build",
			input: "1.2.3",
			valid: false,
		},
		{
			name:  "major minor build revision",
			input: "1.2.3.4",
			valid: true,
		},
		{
			name:  "real-world-example",
			input: "3.6.1000.67",
			valid: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, errors := HDInsightClusterVersion(tt.input, "cluster_version")
			validationFailed := len(errors) > 0

			if tt.valid && validationFailed {
				t.Errorf("Expected %q to be valid but got %+v", tt.input, errors)
			} else if !tt.valid && !validationFailed {
				t.Errorf("Expected %q to be invalid but didn't get an error", tt.input)
			}
		})
	}
}
