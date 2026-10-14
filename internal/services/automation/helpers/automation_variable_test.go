// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"testing"
	"time"
)

func TestParseAzureRmAutomationVariableValue(t *testing.T) {
	type ExpectFunc func(any) bool
	cases := []struct {
		Name        string
		Resource    string
		IsNil       bool
		Value       string
		HasError    bool
		ExpectValue any
		Expect      ExpectFunc
	}{
		{
			Name:        "string variable",
			Resource:    "azurerm_automation_variable_string",
			Value:       "\"Test String\"",
			HasError:    false,
			ExpectValue: "Test String",
			Expect:      func(v any) bool { return v.(string) == "Test String" },
		},
		{
			Name:        "integer variable 135",
			Resource:    "azurerm_automation_variable_int",
			Value:       "135",
			HasError:    false,
			ExpectValue: 135,
			Expect:      func(v any) bool { return v.(int32) == 135 },
		},
		{
			Name:        "integer variable 0",
			Resource:    "azurerm_automation_variable_int",
			Value:       "0",
			HasError:    false,
			ExpectValue: 0,
			Expect:      func(v any) bool { return v.(int32) == 0 },
		},
		{
			Name:        "integer variable 1",
			Resource:    "azurerm_automation_variable_int",
			Value:       "1",
			HasError:    false,
			ExpectValue: 1,
			Expect:      func(v any) bool { return v.(int32) == 1 },
		},
		{
			Name:        "integer variable 2",
			Resource:    "azurerm_automation_variable_int",
			Value:       "2",
			HasError:    false,
			ExpectValue: 2,
			Expect:      func(v any) bool { return v.(int32) == 2 },
		},
		{
			Name:        "boolean variable true",
			Resource:    "azurerm_automation_variable_bool",
			Value:       "true",
			HasError:    false,
			ExpectValue: true,
			Expect:      func(v any) bool { return v.(bool) },
		},
		{
			Name:        "boolean variable false",
			Resource:    "azurerm_automation_variable_bool",
			Value:       "false",
			HasError:    false,
			ExpectValue: false,
			Expect:      func(v any) bool { return !v.(bool) },
		},
		{
			Name:        "datetime variable",
			Resource:    "azurerm_automation_variable_datetime",
			Value:       "\"\\/Date(1556142054074)\\/\"",
			HasError:    false,
			ExpectValue: time.Date(2019, time.April, 24, 21, 40, 54, 74000000, time.UTC),
			Expect: func(v any) bool {
				return v.(time.Time).Equal(time.Date(2019, time.April, 24, 21, 40, 54, 74000000, time.UTC))
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			value := &tc.Value
			if tc.IsNil {
				value = nil
			}
			actual, err := ParseAzureAutomationVariableValue(tc.Resource, value)
			if tc.HasError && err == nil {
				t.Fatalf("Expect parseAzureAutomationVariableValue to return error for resource %q and value %s", tc.Resource, tc.Value)
			}
			if !tc.HasError {
				if err != nil {
					t.Fatalf("Expect parseAzureAutomationVariableValue to return no error for resource %q and value %s, err: %+v", tc.Resource, tc.Value, err)
				} else if !tc.Expect(actual) {
					t.Fatalf("Expect parseAzureAutomationVariableValue to return %v instead of %v for resource %q and value %s", tc.ExpectValue, actual, tc.Resource, tc.Value)
				}
			}
		})
	}
}
