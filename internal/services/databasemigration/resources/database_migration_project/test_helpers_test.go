// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package database_migration_project_test

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

type DatabaseMigrationServiceResource struct{}

func (DatabaseMigrationServiceResource) base(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-dbms-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestVnet-dbms-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestSubnet-dbms-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (DatabaseMigrationServiceResource) basic(data acceptance.TestData) string {
	template := DatabaseMigrationServiceResource{}.base(data)

	return fmt.Sprintf(`
%s

resource "azurerm_database_migration_service" "test" {
  name                = "acctestDbms-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  subnet_id           = azurerm_subnet.test.id
  sku_name            = "Standard_1vCores"
}
`, template, data.RandomInteger)
}
