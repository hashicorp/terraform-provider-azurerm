// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package durabletask_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type SchedulerDataSource struct{}

func TestAccDurableTaskSchedulerDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_durable_task_scheduler", "test")
	d := SchedulerDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("sku_name").HasValue("Consumption"),
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("endpoint").Exists(),
			),
		},
	})
}

func (d SchedulerDataSource) basic(data acceptance.TestData) string {
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

data "azurerm_durable_task_scheduler" "test" {
  name                = azurerm_durable_task_scheduler.test.name
  resource_group_name = azurerm_durable_task_scheduler.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
