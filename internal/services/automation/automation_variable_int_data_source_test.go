// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automation_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type AutomationVariableIntDataSource struct{}

func TestAccDataSourceAzureRMAutomationVariableInt_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_automation_variable_int", "test")
	r := AutomationVariableIntDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("value").HasValue("1234"),
			),
		},
	})
}

func (AutomationVariableIntDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_automation_variable_int" "test" {
  name                    = azurerm_automation_variable_int.test.name
  resource_group_name     = azurerm_automation_variable_int.test.resource_group_name
  automation_account_name = azurerm_automation_variable_int.test.automation_account_name
}
`, AutomationVariableIntResource{}.basic(data))
}
