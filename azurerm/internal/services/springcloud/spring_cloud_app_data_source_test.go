package springcloud_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type SpringCloudAppDataSource struct {
}

func TestAccDataSourceSpringCloudApp_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_spring_cloud_app", "test")
	r := SpringCloudAppDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").Exists(),
			),
		},
	})
}

func (SpringCloudAppDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_spring_cloud_app" "test" {
  name                = azurerm_spring_cloud_app.test.name
  resource_group_name = azurerm_spring_cloud_app.test.resource_group_name
  service_name        = azurerm_spring_cloud_app.test.service_name
}
`, SpringCloudAppResource{}.basic(data))
}
