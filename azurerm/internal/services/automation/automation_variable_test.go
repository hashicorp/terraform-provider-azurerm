package automation_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/automation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
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
			Name:        "string variable",
			Resource:    "azurerm_automation_variable_string",
			Value:       "\"Test String\"",
			HasError:    false,
			ExpectValue: "Test String",
			Expect:      func(v interface{}) bool { return v.(string) == "Test String" },
		},
		{
			Name:        "integer variable 135",
			Resource:    "azurerm_automation_variable_int",
			Value:       "135",
			HasError:    false,
			ExpectValue: 135,
			Expect:      func(v interface{}) bool { return v.(int32) == 135 },
		},
		{
			Name:        "integer variable 0",
			Resource:    "azurerm_automation_variable_int",
			Value:       "0",
			HasError:    false,
			ExpectValue: 0,
			Expect:      func(v interface{}) bool { return v.(int32) == 0 },
		},
		{
			Name:        "integer variable 1",
			Resource:    "azurerm_automation_variable_int",
			Value:       "1",
			HasError:    false,
			ExpectValue: 1,
			Expect:      func(v interface{}) bool { return v.(int32) == 1 },
		},
		{
			Name:        "integer variable 2",
			Resource:    "azurerm_automation_variable_int",
			Value:       "2",
			HasError:    false,
			ExpectValue: 2,
			Expect:      func(v interface{}) bool { return v.(int32) == 2 },
		},
		{
			Name:        "boolean variable true",
			Resource:    "azurerm_automation_variable_bool",
			Value:       "true",
			HasError:    false,
			ExpectValue: true,
			Expect:      func(v interface{}) bool { return v.(bool) == true },
		},
		{
			Name:        "boolean variable false",
			Resource:    "azurerm_automation_variable_bool",
			Value:       "false",
			HasError:    false,
			ExpectValue: false,
			Expect:      func(v interface{}) bool { return v.(bool) == false },
		},
		{
			Name:        "datetime variable",
			Resource:    "azurerm_automation_variable_datetime",
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
			actual, err := automation.ParseAzureAutomationVariableValue(tc.Resource, value)
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

func testCheckAzureRMAutomationVariableExists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState, varType string) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}

	resourceGroup := id.ResourceGroup
	accountName := id.Path["automationAccounts"]
	name := id.Path["variables"]

	resp, err := clients.Automation.VariableClient.Get(ctx, resourceGroup, accountName, name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Automation %s Variable %q (Automation Account Name %q / Resource Group %q) does not exist", varType, name, accountName, resourceGroup)
	}

	return utils.Bool(resp.VariableProperties != nil), nil
}
