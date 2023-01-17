package resourceproviders

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/profiles/2017-03-09/resources/mgmt/resources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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
		enhancedEnabled = features.EnhancedValidationEnabled()
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
	cachedResourceProviders = &[]resources.Provider{{Namespace: utils.String("Microsoft.Compute")}}
	defer func() {
		enhancedEnabled = features.EnhancedValidationEnabled()
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
