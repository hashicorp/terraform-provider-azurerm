package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceAzureRMPublicIP_basic(t *testing.T) {
	dataSourceName := "data.azurerm_public_ip.test"
	ri := acctest.RandInt()

	name := fmt.Sprintf("acctestpublicip-%d", ri)
	resourceGroupName := fmt.Sprintf("acctestRG-%d", ri)

	config := testAccDataSourceAzureRMPublicIPBasic(name, resourceGroupName, ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPublicIpDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", name),
					resource.TestCheckResourceAttr(dataSourceName, "resource_group_name", resourceGroupName),
					resource.TestCheckResourceAttr(dataSourceName, "domain_name_label", fmt.Sprintf("acctest-%d", ri)),
					resource.TestCheckResourceAttr(dataSourceName, "idle_timeout_in_minutes", "30"),
					resource.TestCheckResourceAttrSet(dataSourceName, "fqdn"),
					resource.TestCheckResourceAttrSet(dataSourceName, "ip_address"),
					resource.TestCheckResourceAttr(dataSourceName, "ip_version", "IPv4"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.environment", "test"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMPublicIPBasic(name string, resourceGroupName string, rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "%s"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                         = "%s"
  location                     = "${azurerm_resource_group.test.location}"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  public_ip_address_allocation = "static"
  domain_name_label            = "acctest-%d"
  idle_timeout_in_minutes      = 30

  tags {
    environment = "test"
  }
}

data "azurerm_public_ip" "test" {
  name                = "${azurerm_public_ip.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, resourceGroupName, location, name, rInt)
}
