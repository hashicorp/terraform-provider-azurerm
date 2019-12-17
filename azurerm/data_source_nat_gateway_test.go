package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMNatGateway_basic(t *testing.T) {
	ri := tf.AccRandTimeInt()
	// Using alt location because the resource currently in private preview and is only available in eastus2.
	location := acceptance.AltLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNatGateway_basic(ri, location),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func TestAccDataSourceAzureRMNatGateway_complete(t *testing.T) {
	dataSourceName := "data.azurerm_nat_gateway.test"
	ri := tf.AccRandTimeInt()
	// Using alt location because the resource currently in private preview and is only available in eastus2.
	location := acceptance.AltLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNatGateway_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "public_ip_address_ids.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "public_ip_prefix_ids.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "sku_name", "Standard"),
					resource.TestCheckResourceAttr(dataSourceName, "idle_timeout_in_minutes", "10"),
				),
			},
		},
	})
}

func testAccDataSourceNatGateway_basic(rInt int, location string) string {
	config := testAccAzureRMNatGateway_basic(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_nat_gateway" "test" {
  resource_group_name = "${azurerm_nat_gateway.test.resource_group_name}"
  name                = "${azurerm_nat_gateway.test.name}"
}
`, config)
}

func testAccDataSourceNatGateway_complete(rInt int, location string) string {
	config := testAccAzureRMNatGateway_complete(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_nat_gateway" "test" {
  resource_group_name = "${azurerm_nat_gateway.test.resource_group_name}"
  name                = "${azurerm_nat_gateway.test.name}"
}
`, config)
}
