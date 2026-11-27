// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package system_center_virtual_machine_manager_inventory_items_test

import (
	"fmt"
	"os"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

type SystemCenterVirtualMachineManagerServerResource struct{}

func (r SystemCenterVirtualMachineManagerServerResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_system_center_virtual_machine_manager_server" "test" {
  name                = "acctest-scvmmms-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = "%s"
  fqdn                = "%s"
  username            = "%s"
  password            = "%s"
}
`, r.template(data), data.RandomInteger, os.Getenv("ARM_TEST_CUSTOM_LOCATION_ID"), os.Getenv("ARM_TEST_FQDN"), os.Getenv("ARM_TEST_USERNAME"), os.Getenv("ARM_TEST_PASSWORD"))
}

func (r SystemCenterVirtualMachineManagerServerResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-scvmmms-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}
