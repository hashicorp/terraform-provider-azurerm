// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package systemcentervirtualmachinemanager_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type SystemCenterVirtualMachineManagerInventoryItemsDataSource struct{}

func TestAccDataSourceSystemCenterVirtualMachineManagerInventoryItems_basic(t *testing.T) {
	if os.Getenv("ARM_TEST_CUSTOM_LOCATION_ID") == "" || os.Getenv("ARM_TEST_FQDN") == "" || os.Getenv("ARM_TEST_USERNAME") == "" || os.Getenv("ARM_TEST_PASSWORD") == "" {
		t.Skip("Skipping as one of `ARM_TEST_CUSTOM_LOCATION_ID`, `ARM_TEST_FQDN`, `ARM_TEST_USERNAME`, `ARM_TEST_PASSWORD` was not specified")
	}

	data := acceptance.BuildTestData(t, "data.azurerm_system_center_virtual_machine_manager_inventory_items", "test")
	r := SystemCenterVirtualMachineManagerInventoryItemsDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("inventory_items.#").Exists(),
			),
		},
	})
}

func (d SystemCenterVirtualMachineManagerInventoryItemsDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

// Service team confirmed that waiting to sync Inventory Items in System Center Virtual Machine Manager Server is expected behaviour since the backend operator creates CRD (Custom Resource Definitions) for all the existing resources from onPrem and create InventoryItem resources which takes some time depending upon the number of resources after PUT System Center Virtual Machine Manager Server operation
resource "time_sleep" "wait_1_minute" {
  depends_on = [azurerm_system_center_virtual_machine_manager_server.test]

  create_duration = "1m"
}

data "azurerm_system_center_virtual_machine_manager_inventory_items" "test" {
  inventory_type                                  = "Cloud"
  system_center_virtual_machine_manager_server_id = azurerm_system_center_virtual_machine_manager_server.test.id

  depends_on = [time_sleep.wait_1_minute]
}
`, SystemCenterVirtualMachineManagerServerResource{}.basic(data))
}
