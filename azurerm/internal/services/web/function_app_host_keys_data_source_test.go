package web_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type FunctionAppHostKeysDataSource struct{}

func TestAccFunctionAppHostKeysDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_function_app_host_keys", "test")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: FunctionAppHostKeysDataSource{}.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("primary_key").Exists(),
				check.That(data.ResourceName).Key("default_function_key").Exists(),
				check.That(data.ResourceName).Key("event_grid_extension_config_key").Exists(),
			),
		},
	})
}

func (d FunctionAppHostKeysDataSource) basic(data acceptance.TestData) string {
	template := FunctionAppResource{}.basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_function_app_host_keys" "test" {
  name                = azurerm_function_app.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}
