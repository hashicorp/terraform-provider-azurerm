package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMMsSqlManagedInstance_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_mssql_managed_instance", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlManagedInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMMsSqlManagedInstance_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlManagedInstanceExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttr(data.ResourceName, "location", data.Locations.Primary),
					resource.TestCheckResourceAttr(data.ResourceName, "collation", "SQL_Latin1_General_CP1_CI_AS"),
					resource.TestCheckResourceAttr(data.ResourceName, "proxy_override", "Redirect"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_size_gb", "256"),
					resource.TestCheckResourceAttr(data.ResourceName, "vcores", "16"),
					resource.TestCheckResourceAttr(data.ResourceName, "minimal_tls_version", "1.1"),
					resource.TestCheckResourceAttr(data.ResourceName, "data_endpoint_enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "timezone_id", "UTC"),
					resource.TestCheckResourceAttr(data.ResourceName, "license_type", "LicenseIncluded"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMMsSqlManagedInstance_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_network_security_group" "test" {
  name                = "accTestNetworkSecurityGroup"
  location            = "%[2]s"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-%[1]d-network"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[2]s"
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  virtual_network_name = azurerm_virtual_network.test.name
  resource_group_name  = azurerm_resource_group.test.name
  address_prefixes     = ["10.0.1.0/24"]
  delegation {
    name = "miDelegation"
    service_delegation {
      name = "Microsoft.Sql/managedInstances"
    }
  }
}

resource "azurerm_subnet_network_security_group_association" "test" {
  subnet_id                 = azurerm_subnet.test.id
  network_security_group_id = azurerm_network_security_group.test.id
}

resource "azurerm_route_table" "test" {
  name                = "test-routetable"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  route {
    name                   = "test"
    address_prefix         = "10.100.0.0/14"
    next_hop_type          = "VirtualAppliance"
    next_hop_in_ip_address = "10.10.1.1"
  }
}

resource "azurerm_subnet_route_table_association" "test" {
  subnet_id      = azurerm_subnet.test.id
  route_table_id = azurerm_route_table.test.id
}

resource "azurerm_mssql_managed_instance" "test" {
  name                         = "acctest-mi-%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  administrator_login          = "miadministrator"
  administrator_login_password = "LengthyPassword@1234"
  subnet_id                    = azurerm_subnet.test.id
  identity {
    type = "SystemAssigned"
  }
  sku {
    capacity = 8
    family   = "Gen5"
    name     = "GP_Gen5"
    tier     = "GeneralPurpose"
  }
  license_type          = "LicenseIncluded"
  collation             = "SQL_Latin1_General_CP1_CI_AS"
  proxy_override        = "Redirect"
  storage_size_gb       = 256
  vcores                = 16
  data_endpoint_enabled = true
  timezone_id           = "UTC"
  minimal_tls_version   = "1.1"
}

data "azurerm_mssql_managed_instance" "example" {
  name                = azurerm_mssql_managed_instance.test.name
  resource_group_name = azurerm_mssql_managed_instance.test.resource_group_name

}
`, data.RandomInteger, data.Locations.Primary)
}
