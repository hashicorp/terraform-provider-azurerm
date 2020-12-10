package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMIPGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_ip_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIpGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMIpGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "location"),
					resource.TestCheckResourceAttr(data.ResourceName, "cidrs.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMIpGroup_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_ip_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIpGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMIpGroup_complete(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "location"),
					resource.TestCheckResourceAttr(data.ResourceName, "cidrs.#", "3"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMIpGroup_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_ip_group" "test" {
  name                = azurerm_ip_group.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, testAccAzureRMIpGroup_basic(data))
}

func testAccDataSourceAzureRMIpGroup_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_ip_group" "test" {
  name                = azurerm_ip_group.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, testAccAzureRMIpGroup_complete(data))
}
