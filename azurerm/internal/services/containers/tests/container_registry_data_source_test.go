package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMContainerRegistry_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_container_registry", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMContainerRegistry_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "location"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "admin_enabled"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "login_server"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMContainerRegistry_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_container_registry" "test" {
  name                = azurerm_container_registry.test.name
  resource_group_name = azurerm_container_registry.test.resource_group_name
}
`, testAccAzureRMContainerRegistry_basicManaged(data, "Basic"))
}
