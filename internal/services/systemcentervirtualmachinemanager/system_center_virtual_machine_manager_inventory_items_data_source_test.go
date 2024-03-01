// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package systemcentervirtualmachinemanager_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type SystemCenterVirtualMachineManagerInventoryItemsDataSource struct{}

func TestAccDataSourceSystemCenterVirtualMachineManagerInventoryItems_basic(t *testing.T) {
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
  system_center_virtual_machine_manager_server_id = azurerm_system_center_virtual_machine_manager_server.test.id

  depends_on = [time_sleep.wait_1_minute]
}
`, SystemCenterVirtualMachineManagerServerResource{}.basic(data))
}
