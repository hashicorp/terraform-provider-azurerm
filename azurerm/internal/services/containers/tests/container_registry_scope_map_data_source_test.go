package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMContainerRegistryScopeMap_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_container_registry_scope_map", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMContainerRegistryScopeMap_complete(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "A test scope map"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "container_registry_name"),
					resource.TestCheckResourceAttr(data.ResourceName, "actions.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "actions.0", "repositories/testrepo/content/read"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMContainerRegistryScopeMap_complete(data acceptance.TestData) string {
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
	`, testAccAzureRMContainerRegistry_basicManaged(data, "Premium"), data.RandomInteger)
}
