// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package network_manager_ipam_pool_test

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

type ManagerIpamPoolResource struct{}

func (r ManagerIpamPoolResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_network_manager_ipam_pool" "parent" {
  name               = "acctest-ipampool-p-%[2]d"
  network_manager_id = azurerm_network_manager.test.id
  location           = azurerm_resource_group.test.location
  display_name       = "ipampool0"
  address_prefixes   = ["10.0.0.0/26", "10.0.1.0/27"]
}

resource "azurerm_network_manager_ipam_pool" "test" {
  name               = "acctest-ipampool-%[2]d"
  network_manager_id = azurerm_network_manager.test.id
  location           = azurerm_resource_group.test.location
  display_name       = "ipampool1"
  parent_pool_name   = azurerm_network_manager_ipam_pool.parent.name
  address_prefixes   = ["10.0.0.0/27", "10.0.1.0/28"]
  description        = "This is a test IPAM pool"

  tags = {
    foo = "bar"
    env = "test"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ManagerIpamPoolResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-network-manager-ipampool-%d"
  location = "%s"
}

data "azurerm_subscription" "current" {}

resource "azurerm_network_manager" "test" {
  name                = "acctest-nm-ipam-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  scope {
    subscription_ids = [data.azurerm_subscription.current.id]
  }
  scope_accesses = ["Connectivity"]
}
`, data.RandomInteger, data.Locations.Primary)
}
