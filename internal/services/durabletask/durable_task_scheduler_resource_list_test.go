// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package durabletask_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

func TestAccDurableTaskSchedulerList_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_durable_task_scheduler", "test")

	data.ResourceTest(t, []acceptance.TestStep{
		{
			Config: schedulerListTestConfig_basic(data),
		},
	})
}

func schedulerListTestConfig_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-durabletask-%d"
  location = "%s"
}

resource "azurerm_durable_task_scheduler" "test" {
  name                = "acctestdts%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "Consumption"
  ip_allow_list       = ["0.0.0.0/0"]
}

resource "azurerm_durable_task_scheduler" "test2" {
  name                = "acctestdts%s2"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "Consumption"
  ip_allow_list       = ["0.0.0.0/0"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}
