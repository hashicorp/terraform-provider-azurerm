package validate

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestVirtualNetworkRule_invalidNameValidation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		// Must only contain alphanumeric characters or hyphens (4 cases)
		{
			Value:    "test!Rule",
			ErrCount: 1,
		},
		{
			Value:    "test_Rule",
			ErrCount: 1,
		},
		{
			Value:    "test:Rule",
			ErrCount: 1,
		},
		{
			Value:    "test'Rule",
			ErrCount: 1,
		},
		// Cannot be more than 128 characters (1 case - ensure starts with a letter)
		{
			Value:    fmt.Sprintf("v%s", acctest.RandString(128)),
			ErrCount: 1,
		},
		// Cannot be empty (1 case)
		{
			Value:    "",
			ErrCount: 1,
		},
		// Cannot end in a hyphen (1 case)
		{
			Value:    "testRule-",
			ErrCount: 1,
		},
		// Cannot start with a number or hyphen (2 cases)
		{
			Value:    "7testRule",
			ErrCount: 1,
		},
		{
			Value:    "-testRule",
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := VirtualNetworkRuleName(tc.Value, "azurerm_postgresql_virtual_network_rule")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Bad: Expected the Azure RM PostgreSQL Virtual Network Rule Name to trigger a validation error.")
		}
	}
}

func TestResourceAzureRMPostgreSQLVirtualNetworkRule_validNameValidation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		// Test all lowercase
		{
			Value:    "thisisarule",
			ErrCount: 0,
		},
		// Test all uppercase
		{
			Value:    "THISISARULE",
			ErrCount: 0,
		},
		// Test alternating cases
		{
			Value:    "tHiSiSaRuLe",
			ErrCount: 0,
		},
		// Test hyphens
		{
			Value:    "this-is-a-rule",
			ErrCount: 0,
		},
		// Test multiple hyphens in a row
		{
			Value:    "this----1s----a----ru1e",
			ErrCount: 0,
		},
		// Test numbers (except for first character)
		{
			Value:    "v1108501298509850810258091285091820-5",
			ErrCount: 0,
		},
		// Test a lot of hyphens and numbers
		{
			Value:    "x-5-4-1-2-5-2-6-1-5-2-5-1-2-5-6-2-2",
			ErrCount: 0,
		},
		// Test exactly 128 characters
		{
			Value:    fmt.Sprintf("v%s", acctest.RandString(127)),
			ErrCount: 0,
		},
		// Test short, 1-letter name
		{
			Value:    "V",
			ErrCount: 0,
		},
	}

	for _, tc := range cases {
		_, errors := VirtualNetworkRuleName(tc.Value, "azurerm_postgresql_virtual_network_rule")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Bad: Expected the Virtual Network Rule Name pass name validation successfully but triggered a validation error.")
		}
	}
}
