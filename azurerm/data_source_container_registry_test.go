package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMContainerRegistry_basic(t *testing.T) {
	ri := tf.AccRandTimeInt()
	config := testAccDataSourceAzureRMContainerRegistry_basic(ri)

	dataSourceName := "data.azurerm_container_registry.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "location"),
					resource.TestCheckResourceAttrSet(dataSourceName, "admin_enabled"),
					resource.TestCheckResourceAttrSet(dataSourceName, "login_server"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMContainerRegistry_basic(rInt int) string {
	return fmt.Sprintf(`
%s

data "azurerm_container_registry" "test" {
  name                = "${azurerm_container_registry.test.name}"
  resource_group_name = "${azurerm_container_registry.test.resource_group_name}"
}
`, testAccAzureRMContainerRegistry_basicManaged(rInt, acceptance.Location(), "Basic"))
}
