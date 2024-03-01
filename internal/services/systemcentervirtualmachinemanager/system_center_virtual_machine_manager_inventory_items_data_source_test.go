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

// Inventory Items List API doesn't immediately return the list of Inventory Items after System Center Virtual Machine Manager Server is created
// Once the issue https://github.com/Azure/azure-rest-api-specs/issues/28022 is fixed, this part could be removed
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
