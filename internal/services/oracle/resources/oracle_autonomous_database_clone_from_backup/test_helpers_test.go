// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_autonomous_database_clone_from_backup_test

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

type AdbsRegularResource struct{}

func (a AdbsRegularResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`


%s

provider "azurerm" {
  features {}
}

resource "azurerm_oracle_autonomous_database" "test" {
  name                             = "OFake%[2]d"
  display_name                     = "OFake%[2]d"
  resource_group_name              = azurerm_resource_group.test.name
  location                         = "%[3]s"
  compute_model                    = "ECPU"
  compute_count                    = 2
  license_model                    = "LicenseIncluded"
  backup_retention_period_in_days  = 12
  auto_scaling_enabled             = false
  auto_scaling_for_storage_enabled = false
  mtls_connection_required         = true
  data_storage_size_in_tbs         = 1
  db_workload                      = "APEX"
  admin_password                   = "TestPass#2024#"
  db_version                       = "19c"
  character_set                    = "AL32UTF8"
  national_character_set           = "AL16UTF16"
  subnet_id                        = azurerm_subnet.test.id
  virtual_network_id               = azurerm_virtual_network.test.id
}
`, a.template(data), data.RandomInteger, data.Locations.Primary)
}

func (a AdbsRegularResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest%[1]d_vnet"
  address_space       = ["10.0.0.0/16"]
  location            = "%[2]s"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "eacctest%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]

  delegation {
    name = "delegation"

    service_delegation {
      actions = [
        "Microsoft.Network/networkinterfaces/*",
        "Microsoft.Network/virtualNetworks/subnets/join/action",
      ]
      name = "Oracle.Database/networkAttachments"
    }
  }
}

`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
