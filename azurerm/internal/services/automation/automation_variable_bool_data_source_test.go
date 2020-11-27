package automation_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type AutomationVariableBoolDataSource struct {
}

func TestAccDataSourceAzureRMAutomationVariableBool_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_automation_variable_bool", "test")
	r := AutomationVariableBoolDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("value").HasValue("false"),
			),
		},
	})
}

func (AutomationVariableBoolDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_automation_variable_bool" "test" {
  name                    = azurerm_automation_variable_bool.test.name
  resource_group_name     = azurerm_automation_variable_bool.test.resource_group_name
  automation_account_name = azurerm_automation_variable_bool.test.automation_account_name
}
`, AutomationVariableBoolResource{}.basic(data))
}
