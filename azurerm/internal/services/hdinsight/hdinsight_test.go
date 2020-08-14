package hdinsight

import "testing"

func TestHDInsightClusterVersionDiffSuppress(t *testing.T) {
	tests := []struct {
		name          string
		userInput     string
		azureResponse string
		suppressed    bool
	}{
		{
			name:          "empty name",
			userInput:     "",
			azureResponse: "",
			suppressed:    false,
		},
		{
			name:          "missing user input",
			userInput:     "",
			azureResponse: "1.2.3.4",
			suppressed:    false,
		},
		{
			name:          "missing api response",
			userInput:     "1.2",
			azureResponse: "",
			suppressed:    false,
		},
		{
			name:          "major minor user input",
			userInput:     "3.6",
			azureResponse: "3.6.1000.67",
			suppressed:    true,
		},
		{
			name:          "full version user input",
			userInput:     "3.6.1000.67",
			azureResponse: "3.6.1000.67",
			suppressed:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wasSuppressed := hdinsightClusterVersionDiffSuppressFunc("", tt.userInput, tt.azureResponse, nil)
			if tt.suppressed != wasSuppressed {
				t.Errorf("Expected %q to be %t but got %t", tt.name, tt.suppressed, wasSuppressed)
			}
		})
	}
}
