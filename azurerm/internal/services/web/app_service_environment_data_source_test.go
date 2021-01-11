package web_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type AppServiceEnvironmentDataSource struct{}

func TestAccDataSourceAppServiceEnvironment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_app_service_environment", "test")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: AppServiceEnvironmentDataSource{}.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("front_end_scale_factor").Exists(),
				check.That(data.ResourceName).Key("pricing_tier").Exists(),
				check.That(data.ResourceName).Key("internal_ip_address").Exists(),
				check.That(data.ResourceName).Key("service_ip_address").Exists(),
				check.That(data.ResourceName).Key("outbound_ip_addresses").Exists(),
			),
		},
	})
}

func (d AppServiceEnvironmentDataSource) basic(data acceptance.TestData) string {
	config := AppServiceEnvironmentResource{}.basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_app_service_environment" "test" {
  name                = azurerm_app_service_environment.test.name
  resource_group_name = azurerm_app_service_environment.test.resource_group_name
}
`, config)
}
