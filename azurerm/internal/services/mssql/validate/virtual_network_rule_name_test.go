package validate

import (
	"strings"
	"testing"
)

func TestVirtualNetworkRuleName(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		// Valid cases

		{
			Value:    "thisisarule",
			ErrCount: 0,
		},
		{
			Value:    "THISISARULE",
			ErrCount: 0,
		},
		{
			Value:    "tHiSiSaRuLe",
			ErrCount: 0,
		},
		{
			Value:    "this-is-a-rule",
			ErrCount: 0,
		},
		{
			Value:    "this----1s----a----ru1e",
			ErrCount: 0,
		},
		{
			Value:    "this_is_a_rule",
			ErrCount: 0,
		},
		{
			Value:    "this_is_a_rule_",
			ErrCount: 0,
		},
		{
			Value:    "this____1s____a____ru1e",
			ErrCount: 0,
		},
		{
			Value:    "this.is.a.rule",
			ErrCount: 0,
		},
		{
			Value:    "this....1s....a....ru1e",
			ErrCount: 0,
		},
		{
			Value:    "1108501298509850810258091285091820-5",
			ErrCount: 0,
		},
		{
			Value:    "x-5-4-1-2-5-2-6-1-5-2-5-1-2-5-6-2-2",
			ErrCount: 0,
		},
		{
			Value:    "x_5_4_1_2_5_2_6_1_5_2_5_1_2_5_6_2_2",
			ErrCount: 0,
		},
		{
			Value:    "x.5.4.1.2.5.2.6.1.5.2.5.1.2.5.6.2.2",
			ErrCount: 0,
		},
		{
			Value:    strings.Repeat("a", 64),
			ErrCount: 0,
		},

		// Invalid Cases

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
			Value:    strings.Repeat("a", 65),
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

	for i, tc := range cases {
		_, errors := VirtualNetworkRuleName(tc.Value, "azurerm_sql_virtual_network_rule")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Case %d: Encountered %d error(s), expected %d", i, len(errors), tc.ErrCount)
		}
	}
}
