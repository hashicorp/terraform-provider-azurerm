package validate

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

/*
	--Testing for Failure--
	Validation Function Tests - Invalid Name Validations
*/
func TestVirtualNetworkRuleInvalidNameValidation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		// Must only contain alphanumeric characters, periods, underscores or hyphens (4 cases)
		{
			Value:    "test!Rule",
			ErrCount: 1,
		},
		{
			Value:    "test&Rule",
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
		// Cannot be more than 64 characters (1 case - ensure starts with a letter)
		{
			Value:    fmt.Sprintf("v%s", acceptance.RandString(64)),
			ErrCount: 1,
		},
		// Cannot be empty (1 case)
		{
			Value:    "",
			ErrCount: 1,
		},
		// Cannot be single character (1 case)
		{
			Value:    "a",
			ErrCount: 1,
		},
		// Cannot end in a hyphen (1 case)
		{
			Value:    "testRule-",
			ErrCount: 1,
		},
		// Cannot end in a period (1 case)
		{
			Value:    "testRule.",
			ErrCount: 1,
		},
		// Cannot start with a period, underscore or hyphen (3 cases)
		{
			Value:    ".testRule",
			ErrCount: 1,
		},
		{
			Value:    "_testRule",
			ErrCount: 1,
		},
		{
			Value:    "-testRule",
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := VirtualNetworkRuleName(tc.Value, "azurerm_sql_virtual_network_rule")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Bad: Expected the Azure RM SQL Virtual Network Rule Name to trigger a validation error.")
		}
	}
}

/*
	--Testing for Success--
	Validation Function Tests - (Barely) Valid Name Validations
*/
func TestVirtualNetworkRuleValidNameValidation(t *testing.T) {
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
		// Test underscores
		{
			Value:    "this_is_a_rule",
			ErrCount: 0,
		},
		// Test ending with underscore
		{
			Value:    "this_is_a_rule_",
			ErrCount: 0,
		},
		// Test multiple underscoress in a row
		{
			Value:    "this____1s____a____ru1e",
			ErrCount: 0,
		},
		// Test periods
		{
			Value:    "this.is.a.rule",
			ErrCount: 0,
		},
		// Test multiple periods in a row
		{
			Value:    "this....1s....a....ru1e",
			ErrCount: 0,
		},
		// Test numbers
		{
			Value:    "1108501298509850810258091285091820-5",
			ErrCount: 0,
		},
		// Test a lot of hyphens and numbers
		{
			Value:    "x-5-4-1-2-5-2-6-1-5-2-5-1-2-5-6-2-2",
			ErrCount: 0,
		},
		// Test a lot of underscores and numbers
		{
			Value:    "x_5_4_1_2_5_2_6_1_5_2_5_1_2_5_6_2_2",
			ErrCount: 0,
		},
		// Test a lot of periods and numbers
		{
			Value:    "x.5.4.1.2.5.2.6.1.5.2.5.1.2.5.6.2.2",
			ErrCount: 0,
		},
		// Test exactly 64 characters
		{
			Value:    fmt.Sprintf("v%s", acceptance.RandString(63)),
			ErrCount: 0,
		},
	}

	for _, tc := range cases {
		_, errors := VirtualNetworkRuleName(tc.Value, "azurerm_sql_virtual_network_rule")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Bad: Expected the Azure RM SQL Virtual Network Rule Name pass name validation successfully but triggered a validation error.")
		}
	}
}
