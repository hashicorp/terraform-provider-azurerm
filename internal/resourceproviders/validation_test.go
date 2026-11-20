// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resourceproviders

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
)

func TestEnhancedValidationDisabled(t *testing.T) {
	testCases := []struct {
		input string
		valid bool
	}{
		{
			input: "",
			valid: false,
		},
		{
			input: "micr0soft",
			valid: true,
		},
		{
			input: "microsoft.compute",
			valid: true,
		},
		{
			input: "Microsoft.Compute",
			valid: true,
		},
	}
	enhancedEnabled = false
	defer func() {
		enhancedEnabled = features.EnhancedValidationResourceProvidersEnabled()
		cachedResourceProviders = nil
	}()

	for _, testCase := range testCases {
		t.Logf("Testing %q..", testCase.input)

		warnings, errors := EnhancedValidate(testCase.input, "name")
		valid := len(warnings) == 0 && len(errors) == 0
		if testCase.valid != valid {
			t.Errorf("Expected %t but got %t", testCase.valid, valid)
		}
	}
}

func TestEnhancedValidationEnabled(t *testing.T) {
	testCases := []struct {
		input string
		valid bool
	}{
		{
			input: "",
			valid: false,
		},
		{
			input: "micr0soft",
			valid: false,
		},
		{
			input: "microsoft.compute",
			valid: false,
		},
		{
			input: "Microsoft.Compute",
			valid: true,
		},
	}
	enhancedEnabled = true
	cachedResourceProviders = &[]string{"Microsoft.Compute"}
	defer func() {
		enhancedEnabled = features.EnhancedValidationResourceProvidersEnabled()
		cachedResourceProviders = nil
	}()

	for _, testCase := range testCases {
		t.Logf("Testing %q..", testCase.input)

		warnings, errors := EnhancedValidate(testCase.input, "name")
		valid := len(warnings) == 0 && len(errors) == 0
		if testCase.valid != valid {
			t.Errorf("Expected %t but got %t", testCase.valid, valid)
		}
	}
}

func TestValidateSupported(t *testing.T) {
	testCases := []struct {
		name   string
		cached *[]string
		input  ResourceProviders
		valid  bool
	}{
		{
			name:   "cache unavailable",
			cached: nil,
			input:  ResourceProviders{"Microsoft.Dashboardd": {}},
			valid:  true,
		},
		{
			name:   "supported",
			cached: &[]string{"Microsoft.Compute", "Microsoft.Dashboard"},
			input:  ResourceProviders{"Microsoft.Compute": {}, "Microsoft.Dashboard": {}},
			valid:  true,
		},
		{
			name:   "unsupported",
			cached: &[]string{"Microsoft.Compute", "Microsoft.Dashboard"},
			input:  ResourceProviders{"Microsoft.Compute": {}, "Microsoft.Dashboardd": {}},
			valid:  false,
		},
		{
			name:   "incorrect casing",
			cached: &[]string{"Microsoft.Compute", "Microsoft.Dashboard"},
			input:  ResourceProviders{"microsoft.dashboard": {}},
			valid:  false,
		},
	}
	defer func() {
		cachedResourceProviders = nil
	}()

	for _, testCase := range testCases {
		t.Logf("Testing %q..", testCase.name)

		cachedResourceProviders = testCase.cached
		err := ValidateSupported(testCase.input)
		if valid := err == nil; testCase.valid != valid {
			t.Errorf("Expected %t but got %t: %+v", testCase.valid, valid, err)
		}
	}
}
