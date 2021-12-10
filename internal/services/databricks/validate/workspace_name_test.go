package validate

import (
	"strings"
	"testing"
)

func TestWorkspaceName(t *testing.T) {
	const errEmpty = "cannot be an empty string"
	const errMinLen = "must be at least 3 characters"
	const errMaxLen = "must be no more than 64 characters"
	const errAllowList = "can contain only alphanumeric characters, underscores, and hyphens"

	cases := []struct {
		Name           string
		Input          string
		ExpectedErrors []string
	}{
		// Happy paths:
		{
			Name:  "Entire character allow-list",
			Input: "aZ09_-",
		},
		{
			Name:  "Minimum character length",
			Input: "---",
		},
		{
			Name:  "Maximum character length",
			Input: "0123456789012345678901234567890123456789012345678901234567890123", // 64 chars
		},

		// Simple negative cases:
		{
			Name:           "Introduce a non-allowed character",
			Input:          "aZ09_-$", // dollar sign
			ExpectedErrors: []string{errAllowList},
		},
		{
			Name:           "Below minimum character length",
			Input:          "--",
			ExpectedErrors: []string{errMinLen},
		},
		{
			Name:           "Above maximum character length",
			Input:          "01234567890123456789012345678901234567890123456789012345678901234", // 31 chars
			ExpectedErrors: []string{errMaxLen},
		},
		{
			Name:           "Specifically test for emptiness",
			Input:          "",
			ExpectedErrors: []string{errEmpty},
		},

		// Complex negative cases
		{
			Name:           "Too short and non-allowed char",
			Input:          "*^",
			ExpectedErrors: []string{errMinLen, errAllowList},
		},
		{
			Name:           "Too long and non-allowed char",
			Input:          "0123456789012345678901234567890123456789012345678901234567890123ß",
			ExpectedErrors: []string{errMaxLen, errAllowList},
		},
	}

	errsContain := func(errors []error, text string) bool {
		for _, err := range errors {
			if strings.Contains(err.Error(), text) {
				return true
			}
		}
		return false
	}

	t.Parallel()
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			_, errors := WorkspaceName(tc.Input, "azurerm_databricks_workspace.test.name")

			if len(errors) != len(tc.ExpectedErrors) {
				t.Fatalf("Expected %d errors but got %d for %q: %v", len(tc.ExpectedErrors), len(errors), tc.Input, errors)
			}

			for _, expectedError := range tc.ExpectedErrors {
				if !errsContain(errors, expectedError) {
					t.Fatalf("Errors did not contain expected error: %s", expectedError)
				}
			}
		})
	}
}
