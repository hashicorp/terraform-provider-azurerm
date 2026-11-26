// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_db_nodes_test

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

type CloudVmClusterResource struct{}

func (a CloudVmClusterResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
  %s
resource "azurerm_oracle_cloud_vm_cluster" "test" {
  location                        = "%[3]s"
  name                            = "OFakeVmacctest%[2]d"
  resource_group_name             = azurerm_resource_group.test.name
  cloud_exadata_infrastructure_id = azurerm_oracle_exadata_infrastructure.test.id
  cpu_core_count                  = 4
  data_storage_size_in_tbs        = 2
  db_node_storage_size_in_gbs     = 120
  db_servers                      = [for obj in data.azurerm_oracle_db_servers.test.db_servers : obj.ocid]
  display_name                    = "OFakeVmacctest%[2]d"
  gi_version                      = "23.0.0.0"
  license_model                   = "BringYourOwnLicense"
  memory_size_in_gbs              = 60
  hostname                        = "hostname"
  ssh_public_keys                 = ["ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC+wWK73dCr+jgQOAxNsHAnNNNMEMWOHYEccp6wJm2gotpr9katuF/ZAdou5AaW1C61slRkHRkpRRX9FA9CYBiitZgvCCz+3nWNN7l/Up54Zps/pHWGZLHNJZRYyAB6j5yVLMVHIHriY49d/GZTZVNB8GoJv9Gakwc/fuEZYYl4YDFiGMBP///TzlI4jhiJzjKnEvqPFki5p2ZRJqcbCiF4pJrxUQR/RXqVFQdbRLZgYfJ8xGB878RENq3yQ39d8dVOkq4edbkzwcUmwwwkYVPIoDGsYLaRHnG+To7FvMeyO7xDVQkMKzopTQV8AuKpyvpqu0a9pWOMaiCyDytO7GGN you@me.com"]
  subnet_id                       = azurerm_subnet.virtual_network_subnet.id
  virtual_network_id              = azurerm_virtual_network.virtual_network.id
}`, a.template(data), data.RandomInteger, data.Locations.Primary)
}

func (a CloudVmClusterResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "virtual_network" {
  name                = "OFakeacctest%[1]d_vnet"
  address_space       = ["10.0.0.0/16"]
  location            = "%[2]s"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "virtual_network_subnet" {
  name                 = "OFakeacctest%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.virtual_network.name
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

resource "azurerm_oracle_exadata_infrastructure" "test" {
  name                = "OFakeacctest%[1]d"
  location            = "%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  compute_count       = "2"
  display_name        = "OFakeacctest%[1]d"
  shape               = "Exadata.X9M"
  storage_count       = "3"
  zones               = ["2"]
}

data "azurerm_oracle_db_servers" "test" {
  resource_group_name               = azurerm_resource_group.test.name
  cloud_exadata_infrastructure_name = azurerm_oracle_exadata_infrastructure.test.name
}


`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
