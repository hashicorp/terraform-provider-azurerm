package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMNetworkSecurityGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_network_security_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMNetworkSecurityGroupBasic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "location"),
					resource.TestCheckResourceAttr(data.ResourceName, "security_rule.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMNetworkSecurityGroup_rules(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_network_security_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMNetworkSecurityGroupWithRules(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "location"),
					resource.TestCheckResourceAttr(data.ResourceName, "security_rule.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "security_rule.0.name", "test123"),
					resource.TestCheckResourceAttr(data.ResourceName, "security_rule.0.priority", "100"),
					resource.TestCheckResourceAttr(data.ResourceName, "security_rule.0.direction", "Inbound"),
					resource.TestCheckResourceAttr(data.ResourceName, "security_rule.0.access", "Allow"),
					resource.TestCheckResourceAttr(data.ResourceName, "security_rule.0.protocol", "Tcp"),
					resource.TestCheckResourceAttr(data.ResourceName, "security_rule.0.source_port_range", "*"),
					resource.TestCheckResourceAttr(data.ResourceName, "security_rule.0.destination_port_range", "*"),
					resource.TestCheckResourceAttr(data.ResourceName, "security_rule.0.source_address_prefix", "*"),
					resource.TestCheckResourceAttr(data.ResourceName, "security_rule.0.destination_address_prefix", "*"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMNetworkSecurityGroup_tags(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_network_security_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMNetworkSecurityGroupTags(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "location"),
					resource.TestCheckResourceAttr(data.ResourceName, "security_rule.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "staging"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMNetworkSecurityGroupBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_group" "test" {
  name                = "acctestnsg-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

data "azurerm_network_security_group" "test" {
  name                = azurerm_network_security_group.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccDataSourceAzureRMNetworkSecurityGroupWithRules(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_group" "test" {
  name                = "acctestnsg-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

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
  name                = azurerm_network_security_group.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccDataSourceAzureRMNetworkSecurityGroupTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_group" "test" {
  name                = "acctestnsg-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    environment = "staging"
  }
}

data "azurerm_network_security_group" "test" {
  name                = azurerm_network_security_group.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
