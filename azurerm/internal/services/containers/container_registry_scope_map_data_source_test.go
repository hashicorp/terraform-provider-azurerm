package containers_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type ContainerRegistryScopeMapDataSource struct {
}

func TestAccDataSourceContainerRegistryScopeMap_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_container_registry_scope_map", "test")
	r := ContainerRegistryScopeMapDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("container_registry_name").Exists(),
				check.That(data.ResourceName).Key("description").HasValue("A test scope map"),
				check.That(data.ResourceName).Key("actions.#").HasValue("1"),
				check.That(data.ResourceName).Key("actions.0").HasValue("repositories/testrepo/content/read"),
			),
		},
	})
}

func (ContainerRegistryScopeMapDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_container_registry_scope_map" "test" {
  name                    = "testscopemap%d"
  description             = "A test scope map"
  resource_group_name     = azurerm_resource_group.test.name
  container_registry_name = azurerm_container_registry.test.name
  actions                 = ["repositories/testrepo/content/read"]
}

data "azurerm_container_registry_scope_map" "test" {
  name                    = azurerm_container_registry_scope_map.test.name
  container_registry_name = azurerm_container_registry.test.name
  resource_group_name     = azurerm_container_registry.test.resource_group_name
}
	`, ContainerRegistryResource{}.basicManaged(data, "Premium"), data.RandomInteger)
}
