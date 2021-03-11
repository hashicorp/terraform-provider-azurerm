package automation_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type AutomationVariableIntDataSource struct {
}

func TestAccDataSourceAzureRMAutomationVariableInt_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_automation_variable_int", "test")
	r := AutomationVariableIntDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
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
