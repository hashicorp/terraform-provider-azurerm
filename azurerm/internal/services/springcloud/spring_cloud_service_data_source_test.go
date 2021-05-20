package springcloud_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type SpringCloudServiceDataSource struct {
}

func TestAccDataSourceSpringCloudService_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_spring_cloud_service", "test")
	r := SpringCloudServiceDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("outbound_public_ip_addresses.0").Exists(),
			),
		},
	})
}

func (SpringCloudServiceDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_spring_cloud_service" "test" {
  name                = azurerm_spring_cloud_service.test.name
  resource_group_name = azurerm_spring_cloud_service.test.resource_group_name
}
`, SpringCloudServiceResource{}.basic(data))
}
