package containers_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type ContainerRegistryTokenDataSource struct {
}

func TestAccDataSourceContainerRegistryToken_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_container_registry_token", "test")
	r := ContainerRegistryTokenDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("container_registry_name").Exists(),
				check.That(data.ResourceName).Key("scope_map_id").Exists(),
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
			),
		},
	})
}

func (ContainerRegistryTokenDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_container_registry_scope_map" "pull_repos" {
  name                    = "_repositories_pull"
  container_registry_name = azurerm_container_registry.test.name
  resource_group_name     = azurerm_container_registry.test.resource_group_name
}

resource "azurerm_container_registry_token" "test" {
  name                    = "testtoken%d"
  resource_group_name     = azurerm_resource_group.test.name
  container_registry_name = azurerm_container_registry.test.name
  scope_map_id            = data.azurerm_container_registry_scope_map.pull_repos.id
  enabled                 = true
}

data "azurerm_container_registry_token" "test" {
  name                    = azurerm_container_registry_token.test.name
  container_registry_name = azurerm_container_registry.test.name
  resource_group_name     = azurerm_container_registry.test.resource_group_name
}
	`, ContainerRegistryResource{}.basicManaged(data, "Premium"), data.RandomInteger)
}
