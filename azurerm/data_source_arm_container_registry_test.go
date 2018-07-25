package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceAzureRMContainerRegistry_basic(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccDataSourceAzureRMContainerRegistry_basic(ri)

	acrName := "azurerm_container_registry.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(acrName, "name"),
					resource.TestCheckResourceAttrSet(acrName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(acrName, "location"),
					resource.TestCheckResourceAttrSet(acrName, "admin_enabled"),
					resource.TestCheckResourceAttrSet(acrName, "login_server"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMContainerRegistry_basic(rInt int) string {
	resource := testAccAzureRMContainerRegistry_basicManaged(rInt, testLocation(), "Basic")
	return fmt.Sprintf(`
%s

data "azurerm_container_registry" "test" {
  name                = "${azurerm_container_registry.test.name}"
  resource_group_name = "${azurerm_container_registry.test.resource_group_name}"
}
`, resource)
}
