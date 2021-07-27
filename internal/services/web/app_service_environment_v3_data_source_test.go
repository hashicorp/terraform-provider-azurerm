package web_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type AppServiceEnvironmentV3DataSource struct{}

func TestAccAppServiceEnvironmentV3DataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_app_service_environment_v3", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: AppServiceEnvironmentV3DataSource{}.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("cluster_setting.#").HasValue("2"),
				check.That(data.ResourceName).Key("location").HasValue(data.Locations.Primary),
				check.That(data.ResourceName).Key("ip_ssl_address_count").IsSet(),
				check.That(data.ResourceName).Key("dns_suffix").HasValue(".p.azurewebsites.net"),
				check.That(data.ResourceName).Key("dedicated_host_count").HasValue("0"),
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
