package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMPublicIPPrefix_basic(t *testing.T) {
	ri := tf.AccRandTimeInt()
	name := fmt.Sprintf("acctestpublicipprefix-%d", ri)
	resourceGroupName := fmt.Sprintf("acctestRG-%d", ri)
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMPublicIPPrefixBasic(name, resourceGroupName, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.azurerm_public_ip_prefix.test", "name", name),
					resource.TestCheckResourceAttr("data.azurerm_public_ip_prefix.test", "resource_group_name", resourceGroupName),
					resource.TestCheckResourceAttr("data.azurerm_public_ip_prefix.test", "location", location),
					resource.TestCheckResourceAttr("data.azurerm_public_ip_prefix.test", "sku", "Standard"),
					resource.TestCheckResourceAttr("data.azurerm_public_ip_prefix.test", "prefix_length", "31"),
					resource.TestCheckResourceAttrSet("data.azurerm_public_ip_prefix.test", "ip_prefix"),
					resource.TestCheckResourceAttr("data.azurerm_public_ip_prefix.test", "tags.%", "1"),
					resource.TestCheckResourceAttr("data.azurerm_public_ip_prefix.test", "tags.env", "test"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMPublicIPPrefixBasic(name string, resourceGroupName string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "%s"
  location = "%s"

  tags = {
    env = "test"
  }
}

resource "azurerm_public_ip_prefix" "test" {
  name                = "%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  prefix_length       = 31

  tags = {
    env = "test"
  }
}

data "azurerm_public_ip_prefix" "test" {
  name                = azurerm_public_ip_prefix.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, resourceGroupName, location, name)
}
