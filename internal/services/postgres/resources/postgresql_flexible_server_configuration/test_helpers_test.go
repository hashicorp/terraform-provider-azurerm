// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package postgresql_flexible_server_configuration_test

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

type PostgresqlFlexibleServerResource struct{}

func (r PostgresqlFlexibleServerResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_flexible_server" "test" {
  name                   = "acctest-fs-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "adminTerraform"
  administrator_password = "QAZwsx123"
  version                = "16"
  sku_name               = "GP_Standard_D2s_v3"
  zone                   = "2"
}
`, r.template(data), data.RandomInteger)
}

func (r PostgresqlFlexibleServerResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vn-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-sn-%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
  service_endpoint {
    service = "Microsoft.Storage"
  }
  delegation {
    name = "fs"
    service_delegation {
      name = "Microsoft.DBforPostgreSQL/flexibleServers"
      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
      ]
    }
  }
}

resource "azurerm_private_dns_zone" "test" {
  name                = "acc%[2]d.postgres.database.azure.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_private_dns_zone_virtual_network_link" "test" {
  name                = "acctestVnetZone%[2]d.com"
  private_dns_zone_id = azurerm_private_dns_zone.test.id
  virtual_network_id  = azurerm_virtual_network.test.id

  depends_on = [azurerm_subnet.test]
}

resource "azurerm_postgresql_flexible_server" "test" {
  name                          = "acctest-fs-%[2]d"
  resource_group_name           = azurerm_resource_group.test.name
  location                      = azurerm_resource_group.test.location
  administrator_login           = "adminTerraform"
  administrator_password        = "QAZwsx123"
  version                       = "16"
  backup_retention_days         = 7
  storage_mb                    = 32768
  delegated_subnet_id           = azurerm_subnet.test.id
  private_dns_zone_id           = azurerm_private_dns_zone.test.id
  public_network_access_enabled = false
  sku_name                      = "GP_Standard_D2s_v3"
  zone                          = "1"

  high_availability {
    mode                      = "ZoneRedundant"
    standby_availability_zone = "2"
  }

  maintenance_window {
    day_of_week  = 0
    start_hour   = 8
    start_minute = 0
  }

  tags = {
    ENV = "Test"
  }

  depends_on = [azurerm_private_dns_zone_virtual_network_link.test]
}
`, r.template(data), data.RandomInteger)
}

func (PostgresqlFlexibleServerResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-postgresql-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Secondary)
}
