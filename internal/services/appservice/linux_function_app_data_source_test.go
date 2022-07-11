package appservice_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type LinuxFunctionAppDataSource struct{}

func TestAccLinuxFunctionAppDataSource_standardComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_linux_function_app", "test")
	d := LinuxFunctionAppDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.standardComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").HasValue(data.Locations.Primary),
				check.That(data.ResourceName).Key("default_hostname").HasValue(fmt.Sprintf("acctest-lfa-%d.azurewebsites.net", data.RandomInteger)),
			),
		},
	})
}

func (LinuxFunctionAppDataSource) standardComplete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_linux_function_app" "test" {
  name                = azurerm_linux_function_app.test.name
  resource_group_name = azurerm_linux_function_app.test.resource_group_name
}
`, LinuxFunctionAppResource{}.standardComplete(data))
}
