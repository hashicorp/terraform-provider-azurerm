// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestSkuProfileVMSizeName(t *testing.T) {
	testCases := []struct {
		name        string
		input       interface{}
		shouldError bool
	}{
		{
			name:        "empty string",
			input:       "",
			shouldError: true,
		},
		{
			name:        "missing Standard prefix",
			input:       "D2s_v5",
			shouldError: true,
		},
		{
			name:        "supported A family",
			input:       "Standard_A2_v2",
			shouldError: false,
		},
		{
			name:        "supported B family",
			input:       "Standard_B2s",
			shouldError: false,
		},
		{
			name:        "supported D family",
			input:       "Standard_D2s_v5",
			shouldError: false,
		},
		{
			name:        "supported E family",
			input:       "Standard_E2s_v5",
			shouldError: false,
		},
		{
			name:        "supported F family",
			input:       "Standard_F2s_v2",
			shouldError: false,
		},
		{
			name:        "unsupported L family",
			input:       "Standard_L8s_v3",
			shouldError: true,
		},
		{
			name:        "unsupported DC family",
			input:       "Standard_DC2ads_v5",
			shouldError: true,
		},
		{
			name:        "unsupported EC family",
			input:       "Standard_EC2as_v5",
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
			_, errors := SkuProfileVMSizeName(testCase.input, "test")

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
