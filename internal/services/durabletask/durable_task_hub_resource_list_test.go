// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package durabletask_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

func TestAccDurableTaskHubList_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_durable_task_task_hub", "test")
	r := TaskHubResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: taskHubListTestConfig_basic(data),
		},
	})
}

func taskHubListTestConfig_basic(data acceptance.TestData) string {
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

resource "azurerm_durable_task_task_hub" "test" {
  name                = "acctest%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  scheduler_id        = azurerm_durable_task_scheduler.test.id
}

resource "azurerm_durable_task_task_hub" "test2" {
  name                = "acctest%s2"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  scheduler_id        = azurerm_durable_task_scheduler.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomString)
}
