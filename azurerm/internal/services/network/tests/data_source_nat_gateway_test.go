package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMNatGateway_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_nat_gateway", "test")
	// Using alt location because the resource currently in private preview and is only available in eastus2.

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNatGateway_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "location"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMNatGateway_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_nat_gateway", "test")
	// Using alt location because the resource currently in private preview and is only available in eastus2.

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNatGateway_complete(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "public_ip_address_ids.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "public_ip_prefix_ids.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "Standard"),
					resource.TestCheckResourceAttr(data.ResourceName, "idle_timeout_in_minutes", "10"),
				),
			},
		},
	})
}

func testAccDataSourceNatGateway_basic(data acceptance.TestData) string {
	config := testAccAzureRMNatGateway_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_nat_gateway" "test" {
  resource_group_name = azurerm_nat_gateway.test.resource_group_name
  name                = azurerm_nat_gateway.test.name
}
`, config)
}

func testAccDataSourceNatGateway_complete(data acceptance.TestData) string {
	config := testAccAzureRMNatGateway_complete(data)
	return fmt.Sprintf(`
%s

data "azurerm_nat_gateway" "test" {
  resource_group_name = azurerm_nat_gateway.test.resource_group_name
  name                = azurerm_nat_gateway.test.name
}
`, config)
}
