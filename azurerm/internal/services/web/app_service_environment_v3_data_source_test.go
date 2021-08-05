package web_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type AppServiceEnvironmentV3DataSource struct{}

func TestAccAppServiceEnvironmentV3DataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_app_service_environment_v3", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: AppServiceEnvironmentV3DataSource{}.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("cluster_setting.#").HasValue("2"),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
			),
		},
	})
}

func (AppServiceEnvironmentV3DataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_app_service_environment_v3" "test" {
  name                = azurerm_app_service_environment_v3.test.name
  resource_group_name = azurerm_app_service_environment_v3.test.resource_group_name
}
`, AppServiceEnvironmentV3Resource{}.complete(data))
}
