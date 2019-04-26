package azure

import (
	"testing"
	"time"
)

func TestParseAzureRmAutomationVariableValue(t *testing.T) {
	type ExpectFunc func(interface{}) bool
	cases := []struct {
		Name        string
		Resource    string
		IsNil       bool
		Value       string
		HasError    bool
		ExpectValue interface{}
		Expect      ExpectFunc
	}{
		{
			Name:        "null variable",
			Resource:    "azurerm_automation_null_variable",
			IsNil:       true,
			Value:       "<nil>",
			HasError:    false,
			ExpectValue: nil,
			Expect:      func(v interface{}) bool { return v == nil },
		},
		{
			Name:        "string variable",
			Resource:    "azurerm_automation_string_variable",
			Value:       "\"Test String\"",
			HasError:    false,
			ExpectValue: "Test String",
			Expect:      func(v interface{}) bool { return v.(string) == "Test String" },
		},
		{
			Name:        "integer variable",
			Resource:    "azurerm_automation_int_variable",
			Value:       "135",
			HasError:    false,
			ExpectValue: 135,
			Expect:      func(v interface{}) bool { return v.(int32) == 135 },
		},
		{
			Name:        "boolean variable",
			Resource:    "azurerm_automation_bool_variable",
			Value:       "true",
			HasError:    false,
			ExpectValue: true,
			Expect:      func(v interface{}) bool { return v.(bool) == true },
		},
		{
			Name:        "datetime variable",
			Resource:    "azurerm_automation_datetime_variable",
			Value:       "\"\\/Date(1556142054074)\\/\"",
			HasError:    false,
			ExpectValue: time.Date(2019, time.April, 24, 21, 40, 54, 74000000, time.UTC),
			Expect: func(v interface{}) bool {
				return v.(time.Time) == time.Date(2019, time.April, 24, 21, 40, 54, 74000000, time.UTC)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			value := &tc.Value
			if tc.IsNil {
				value = nil
			}
			actual, err := ParseAzureRmAutomationVariableValue(tc.Resource, value)
			if tc.HasError && err == nil {
				t.Fatalf("Expect parseAzureRmAutomationVariableValue to return error for resource %q and value %s", tc.Resource, tc.Value)
			}
			if !tc.HasError {
				if err != nil {
					t.Fatalf("Expect parseAzureRmAutomationVariableValue to return no error for resource %q and value %s, err: %+v", tc.Resource, tc.Value, err)
				} else if !tc.Expect(actual) {
					t.Fatalf("Expect parseAzureRmAutomationVariableValue to return %v instead of %v for resource %q and value %s", tc.ExpectValue, actual, tc.Resource, tc.Value)
				}
			}
		})
	}
}
