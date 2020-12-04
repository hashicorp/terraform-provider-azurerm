package validate

import "testing"

func TestTemplateDeploymentContentVersion(t *testing.T) {
	testCases := []struct {
		input string
		valid bool
	}{
		{input: "", valid: false},
		{input: "1.0.0", valid: false},
		{input: "1.0.0.0.0", valid: false},
		{input: "1.0.0.0", valid: true},
		{input: "12.345.24.12", valid: true},
	}

	for _, testCase := range testCases {
		t.Logf("Testing %q..", testCase.input)
		warnings, errors := TemplateDeploymentContentVersion(testCase.input, "content_version")
		valid := len(warnings) == 0 && len(errors) == 0

		if valid != testCase.valid {
			t.Fatalf("Expected %t but got %t - %d warnings %d errors", testCase.valid, valid, len(warnings), len(errors))
		}
	}
}
