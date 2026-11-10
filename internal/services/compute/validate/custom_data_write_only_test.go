// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"strings"
	"testing"
)

func TestCustomDataWriteOnly(t *testing.T) {
	testCases := []struct {
		name        string
		input       any
		shouldError bool
	}{
		{
			name:        "base64 encoded",
			input:       "IyEvYmluL2Jhc2g=",
			shouldError: false,
		},
		{
			name:        "empty string",
			input:       "",
			shouldError: true,
		},
		{
			name:        "not base64 encoded",
			input:       "#!/bin/bash",
			shouldError: true,
		},
		{
			name:        "base64 encoded without padding",
			input:       "IyEvYmluL2Jhc2g",
			shouldError: true,
		},
		{
			name:        "non-string input",
			input:       1,
			shouldError: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, errors := CustomDataWriteOnly(testCase.input, "custom_data_wo")

			hasErrors := len(errors) > 0
			if testCase.shouldError && !hasErrors {
				t.Fatalf("expected an error but got none")
			}

			if !testCase.shouldError && hasErrors {
				t.Fatalf("expected no errors but got %d", len(errors))
			}
		})
	}
}

func TestCustomDataWriteOnly_errorOmitsValue(t *testing.T) {
	input := "#!/bin/bash"

	_, errors := CustomDataWriteOnly(input, "custom_data_wo")
	if len(errors) == 0 {
		t.Fatalf("expected an error but got none")
	}

	for _, err := range errors {
		if strings.Contains(err.Error(), input) {
			t.Fatalf("expected the error to omit the value but got: %s", err)
		}
	}
}
