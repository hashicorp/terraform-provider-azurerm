package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMContainerRegistryToken_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_container_registry_token", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMContainerRegistryToken_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "container_registry_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "scope_map_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "status", "disabled"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMContainerRegistryToken_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_container_registry_scope_map" "pull_repos" {
	name = "_repositories_pull"
	container_registry_name = azurerm_container_registry.test.name
	resource_group_name = azurerm_container_registry.test.resource_group_name
}

resource "azurerm_container_registry_token" "test" {
	name                = "testtoken%d"
	resource_group_name = azurerm_resource_group.test.name
	container_registry_name = azurerm_container_registry.test.name
	scope_map_id = data.azurerm_container_registry_scope_map.pull_repos.id
	status = "disabled"
}

data "azurerm_container_registry_token" "test" {
	name = azurerm_container_registry_token.test.name
	container_registry_name = azurerm_container_registry.test.name
	resource_group_name = azurerm_container_registry.test.resource_group_name
}
	`, testAccAzureRMContainerRegistry_basicManaged(data, "Premium"), data.RandomInteger)
}
