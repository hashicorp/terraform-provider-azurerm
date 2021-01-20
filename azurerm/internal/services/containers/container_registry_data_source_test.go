package containers_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type ContainerRegistryDataSource struct {
}

func TestAccDataSourceAzureRMContainerRegistry_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_container_registry", "test")
	r := ContainerRegistryDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("admin_enabled").Exists(),
				check.That(data.ResourceName).Key("login_server").Exists(),
			),
		},
	})
}

func (ContainerRegistryDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_container_registry" "test" {
  name                = azurerm_container_registry.test.name
  resource_group_name = azurerm_container_registry.test.resource_group_name
}
`, ContainerRegistryResource{}.basicManaged(data, "Basic"))
}
