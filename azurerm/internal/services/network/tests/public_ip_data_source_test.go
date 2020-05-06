package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMPublicIP_static(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_public_ip", "test")

	name := fmt.Sprintf("acctestpublicip-%d", data.RandomInteger)
	resourceGroupName := fmt.Sprintf("acctestRG-%d", data.RandomInteger)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPublicIpDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMPublicIP_static(name, resourceGroupName, data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "name", name),
					resource.TestCheckResourceAttr(data.ResourceName, "resource_group_name", resourceGroupName),
					resource.TestCheckResourceAttr(data.ResourceName, "domain_name_label", fmt.Sprintf("acctest-%d", data.RandomInteger)),
					resource.TestCheckResourceAttr(data.ResourceName, "idle_timeout_in_minutes", "30"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "fqdn"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ip_address"),
					resource.TestCheckResourceAttr(data.ResourceName, "ip_version", "IPv4"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "test"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMPublicIP_dynamic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_public_ip", "test")

	name := fmt.Sprintf("acctestpublicip-%d", data.RandomInteger)
	resourceGroupName := fmt.Sprintf("acctestRG-%d", data.RandomInteger)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPublicIpDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMPublicIP_dynamic(data, "Ipv4"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "name", name),
					resource.TestCheckResourceAttr(data.ResourceName, "resource_group_name", resourceGroupName),
					resource.TestCheckResourceAttr(data.ResourceName, "domain_name_label", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "fqdn", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "ip_address", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "ip_version", "IPv4"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "test"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMPublicIP_static(name string, resourceGroupName string, data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "%s"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                    = "%s"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  allocation_method       = "Static"
  domain_name_label       = "acctest-%d"
  idle_timeout_in_minutes = 30

  tags = {
    environment = "test"
  }
}

data "azurerm_public_ip" "test" {
  name                = azurerm_public_ip.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, resourceGroupName, data.Locations.Primary, name, data.RandomInteger)
}

func testAccDataSourceAzureRMPublicIP_dynamic(data acceptance.TestData, ipVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"

  ip_version = "%s"

  tags = {
    environment = "test"
  }
}

data "azurerm_public_ip" "test" {
  name                = azurerm_public_ip.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, ipVersion)
}
