package location

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
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
			input: "chinanorth",
			valid: true,
		},
		{
			input: "China North",
			valid: true,
		},
		{
			input: "westeurope",
			valid: true,
		},
		{
			input: "West Europe",
			valid: true,
		},
		{
			input: "global",
			valid: true,
		},
	}
	enhancedEnabled = false
	defer func() {
		enhancedEnabled = features.EnhancedValidationEnabled()
	}()

	for _, testCase := range testCases {
		t.Logf("Testing %q..", testCase.input)

		warnings, errors := EnhancedValidate(testCase.input, "location")
		valid := len(warnings) == 0 && len(errors) == 0
		if testCase.valid != valid {
			t.Errorf("Expected %t but got %t", testCase.valid, valid)
		}
	}
}

func TestEnhancedValidationEnabledButIsOffline(t *testing.T) {
	testCases := []struct {
		input string
		valid bool
	}{
		{
			input: "",
			valid: false,
		},
		{
			input: "chinanorth",
			valid: true,
		},
		{
			input: "China North",
			valid: true,
		},
		{
			input: "westeurope",
			valid: true,
		},
		{
			input: "West Europe",
			valid: true,
		},
		{
			input: "global",
			valid: true,
		},
	}
	enhancedEnabled = true
	supportedLocations = nil
	defer func() {
		enhancedEnabled = features.EnhancedValidationEnabled()
	}()

	for _, testCase := range testCases {
		t.Logf("Testing %q..", testCase.input)

		warnings, errors := EnhancedValidate(testCase.input, "location")
		valid := len(warnings) == 0 && len(errors) == 0
		if testCase.valid != valid {
			t.Logf("Expected %t but got %t", testCase.valid, valid)
			t.Fail()
		}
	}
}

func TestEnhancedValidationEnabled(t *testing.T) {
	testCases := []struct {
		availableLocations []string
		input              string
		valid              bool
	}{
		{
			availableLocations: publicLocations,
			input:              "",
			valid:              false,
		},
		{
			availableLocations: publicLocations,
			input:              "chinanorth",
			valid:              false,
		},
		{
			availableLocations: publicLocations,
			input:              "China North",
			valid:              false,
		},
		{
			availableLocations: publicLocations,
			input:              "westeurope",
			valid:              true,
		},
		{
			availableLocations: publicLocations,
			input:              "West Europe",
			valid:              true,
		},
		{
			availableLocations: chinaLocations,
			input:              "chinanorth",
			valid:              true,
		},
		{
			availableLocations: chinaLocations,
			input:              "China North",
			valid:              true,
		},
		{
			availableLocations: chinaLocations,
			input:              "westeurope",
			valid:              false,
		},
		{
			availableLocations: chinaLocations,
			input:              "West Europe",
			valid:              false,
		},
		{
			availableLocations: publicLocations,
			input:              "global",
			valid:              true,
		},
	}
	enhancedEnabled = true
	defer func() {
		enhancedEnabled = features.EnhancedValidationEnabled()
		supportedLocations = nil
	}()

	for _, testCase := range testCases {
		t.Logf("Testing %q..", testCase.input)
		supportedLocations = &testCase.availableLocations

		warnings, errors := EnhancedValidate(testCase.input, "location")
		valid := len(warnings) == 0 && len(errors) == 0
		if testCase.valid != valid {
			t.Logf("Expected %t but got %t", testCase.valid, valid)
			t.Fail()
		}
	}
}

var chinaLocations = []string{"chinaeast", "chinanorth", "chinanorth2", "chinaeast2"}
var publicLocations = []string{
	"westus",
	"westus2",
	"eastus",
	"centralus",
	"southcentralus",
	"northcentralus",
	"westcentralus",
	"eastus2",
	"brazilsouth",
	"brazilus",
	"northeurope",
	"westeurope",
	"eastasia",
	"southeastasia",
	"japanwest",
	"japaneast",
	"koreacentral",
	"koreasouth",
	"indiasouth",
	"indiawest",
	"indiacentral",
	"australiaeast",
	"australiasoutheast",
	"canadacentral",
	"canadaeast",
	"uknorth",
	"uksouth2",
	"uksouth",
	"ukwest",
	"francecentral",
	"francesouth",
	"australiacentral",
	"australiacentral2",
	"uaecentral",
	"uaenorth",
	"southafricanorth",
	"southafricawest",
	"switzerlandnorth",
	"switzerlandwest",
	"germanynorth",
	"germanywestcentral",
	"norwayeast",
	"norwaywest",
}
