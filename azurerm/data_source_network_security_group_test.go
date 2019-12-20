package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMNetworkSecurityGroup_basic(t *testing.T) {
	dataSourceName := "data.azurerm_network_security_group.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	config := testAccDataSourceAzureRMNetworkSecurityGroupBasic(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "location"),
					resource.TestCheckResourceAttr(dataSourceName, "security_rule.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMNetworkSecurityGroup_rules(t *testing.T) {
	dataSourceName := "data.azurerm_network_security_group.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	config := testAccDataSourceAzureRMNetworkSecurityGroupWithRules(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "location"),
					resource.TestCheckResourceAttr(dataSourceName, "security_rule.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "security_rule.0.name", "test123"),
					resource.TestCheckResourceAttr(dataSourceName, "security_rule.0.priority", "100"),
					resource.TestCheckResourceAttr(dataSourceName, "security_rule.0.direction", "Inbound"),
					resource.TestCheckResourceAttr(dataSourceName, "security_rule.0.access", "Allow"),
					resource.TestCheckResourceAttr(dataSourceName, "security_rule.0.protocol", "Tcp"),
					resource.TestCheckResourceAttr(dataSourceName, "security_rule.0.source_port_range", "*"),
					resource.TestCheckResourceAttr(dataSourceName, "security_rule.0.destination_port_range", "*"),
					resource.TestCheckResourceAttr(dataSourceName, "security_rule.0.source_address_prefix", "*"),
					resource.TestCheckResourceAttr(dataSourceName, "security_rule.0.destination_address_prefix", "*"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMNetworkSecurityGroup_tags(t *testing.T) {
	dataSourceName := "data.azurerm_network_security_group.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	config := testAccDataSourceAzureRMNetworkSecurityGroupTags(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "location"),
					resource.TestCheckResourceAttr(dataSourceName, "security_rule.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.environment", "staging"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMNetworkSecurityGroupBasic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_group" "test" {
  name                = "acctestnsg-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

data "azurerm_network_security_group" "test" {
  name                = "${azurerm_network_security_group.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, rInt, location, rInt)
}

func testAccDataSourceAzureRMNetworkSecurityGroupWithRules(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_group" "test" {
  name                = "acctestnsg-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  security_rule {
    name                       = "test123"
    priority                   = 100
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "*"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }
}

data "azurerm_network_security_group" "test" {
  name                = "${azurerm_network_security_group.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, rInt, location, rInt)
}

func testAccDataSourceAzureRMNetworkSecurityGroupTags(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_group" "test" {
  name                = "acctestnsg-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  tags = {
    environment = "staging"
  }
}

data "azurerm_network_security_group" "test" {
  name                = "${azurerm_network_security_group.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, rInt, location, rInt)
}
