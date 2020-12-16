package automation_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type AutomationVariableStringDataSource struct {
}

func TestAccDataSourceAzureRMAutomationVariableString_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_automation_variable_string", "test")
	r := AutomationVariableStringDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("value").Exists(),
			),
		},
	})
}

func (AutomationVariableStringDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_automation_variable_string" "test" {
  name                    = azurerm_automation_variable_string.test.name
  resource_group_name     = azurerm_automation_variable_string.test.resource_group_name
  automation_account_name = azurerm_automation_variable_string.test.automation_account_name
}
`, AutomationVariableStringResource{}.basic(data))
}
