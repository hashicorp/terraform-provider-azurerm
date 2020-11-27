package automation_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type AutomationVariableDateTimeDataSouce struct {
}

func TestAccDataSourceAzureRMAutomationVariableDateTime_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_automation_variable_datetime", "test")
	r := AutomationVariableDateTimeDataSouce{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("value").HasValue("2019-04-24T21:40:54.074Z"),
			),
		},
	})
}

func (AutomationVariableDateTimeDataSouce) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_automation_variable_datetime" "test" {
  name                    = azurerm_automation_variable_datetime.test.name
  resource_group_name     = azurerm_automation_variable_datetime.test.resource_group_name
  automation_account_name = azurerm_automation_variable_datetime.test.automation_account_name
}
`, AutomationVariableDateTimeResource{}.basic(data))
}
